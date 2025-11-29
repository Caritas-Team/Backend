package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/Caritas-Team/reviewer/internal/config"
	"github.com/Caritas-Team/reviewer/internal/handler"
	"github.com/Caritas-Team/reviewer/internal/logger"
	"github.com/Caritas-Team/reviewer/internal/memcached"
	"github.com/Caritas-Team/reviewer/internal/metrics"
	"github.com/Caritas-Team/reviewer/internal/usecase/file"
	"github.com/Caritas-Team/reviewer/internal/usecase/user"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func main() {
	// Конфиг
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load error", "err", err)
		return
	}

	// Базовый контекст
	ctx := context.Background()

	// Глобальный логер
	logger.InitGlobalLogger(cfg)

	// Jaeger
	shutdownTracer, err := initTracer("reviewer", cfg.Jaeger.Endpoint)
	if err != nil {
		slog.Error("tracer init error", "err", err)
		return
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracer(ctx); err != nil {
			slog.Error("tracer shutdown error", "err", err)
		}
	}()

	// Контекст, отменяемый по SIGINT/SIGTERM
	rootCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Кэш
	cache, err := memcached.NewCache(rootCtx, cfg)
	if err != nil {
		slog.Error("cache initialization failed", "err", err)
		return
	}

	rateLimiter := user.NewRateLimiter(cache, cfg)
	rateLimiterMiddleware := handler.NewRateLimiterMiddleware(rateLimiter)

	fileCleaner := file.NewFileCleaner(cache)

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			if err := fileCleaner.DeleteDownloadedFiles(ctx); err != nil {
				slog.Error("file cleaner delete error", "err", err)
			} else {
				slog.Info("file cleaner deleted successfully")
			}
		}
	}()

	// ready + HTTP маршруты для тестов и прочего
	var ready atomic.Bool
	ready.Store(true)

	mux := http.NewServeMux()

	// Для теста CORS. МОЖНО УДАЛЯТЬ
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})

	// Для проверки, готов ли сервер принимать новый трафик. МОЖНО УДАЛЯТЬ
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if ready.Load() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
			return
		}
		http.Error(w, "shutting down", http.StatusServiceUnavailable)
	})

	// Метрики
	metrics.InitMetricsOn(mux)

	// CORS
	h := handler.CORS(handler.CORSConfig{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   cfg.CORS.AllowedMethods,
		AllowedHeaders:   cfg.CORS.AllowedHeaders,
		AllowCredentials: true,
		MaxAgeSeconds:    3600,
	})(mux)

	h = rateLimiterMiddleware.Handler(h)

	h = otelhttp.NewHandler(h, "http-server")

	// HTTP сервер
	srv := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      h,
		ReadTimeout:  cfg.Server.ReadTimeout(),
		WriteTimeout: cfg.Server.WriteTimeout(),
		IdleTimeout:  5 * time.Minute,
	}

	// Запуск сервера
	errCh := make(chan error, 1)
	go func() {
		slog.Info("http server listening", "addr", srv.Addr, "pid", os.Getpid())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	// Сигнал или ошибка сервера
	select {
	case <-rootCtx.Done():
		slog.Info("shutdown signal received")
	case err := <-errCh:
		if err != nil {
			slog.Error("server failed", "err", err)
			_ = cache.Close()
			return
		}
	}

	// Graceful shutdown
	ready.Store(false)

	graceTime := 30 * time.Second
	shCtx, cancel := context.WithTimeout(ctx, graceTime)
	defer cancel()

	if err := srv.Shutdown(shCtx); err != nil {
		slog.Error("http shutdown error", "err", err)
	} else {
		slog.Info("http server shutdown complete")
	}

	if err := cache.Close(); err != nil {
		slog.Error("cache close error", "err", err)
	} else {
		slog.Info("cache closed")
	}

	slog.Info("graceful shutdown finished")

}

func initTracer(serviceName string, endpoint string) (func(context.Context) error, error) {
	if endpoint == "" {
		endpoint = "http://jaeger:14268/api/traces"
	}

	exp, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(endpoint),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(serviceName),
			),
		),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp.Shutdown, nil
}
