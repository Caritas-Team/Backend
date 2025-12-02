package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Caritas-Team/reviewer/internal/check"
	"github.com/Caritas-Team/reviewer/internal/config"
	"github.com/Caritas-Team/reviewer/internal/handler"
	"github.com/Caritas-Team/reviewer/internal/logger"
	"github.com/Caritas-Team/reviewer/internal/memcached"
	"github.com/Caritas-Team/reviewer/internal/metrics"
	"github.com/Caritas-Team/reviewer/internal/usecase/file"
	"github.com/Caritas-Team/reviewer/internal/usecase/user"
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
	log := logger.NewLogger(cfg)

	// Контекст, отменяемый по SIGINT/SIGTERM
	rootCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Кэш
	cache, err := memcached.NewCache(rootCtx, cfg)
	if err != nil {
		log.Error("cache initialization failed", "err", err)
		return
	}

	rateLimiter := user.NewRateLimiter(cache, cfg)
	rateLimiterMiddleware := handler.NewRateLimiterMiddleware(rateLimiter)

	fileCleaner := file.NewFileCleaner(log, cache)

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			if err := fileCleaner.DeleteDownloadedFiles(ctx); err != nil {
				log.Error("file cleaner delete error", "err", err)
			} else {
				log.Info("file cleaner deleted successfully")
			}
		}
	}()

	// Создаём экземпляр ReadinessChecker
	checker := check.NewReadinessChecker(cache, rateLimiterMiddleware)

	mux := http.NewServeMux()

	// Эндпоинт для health check
	mux.HandleFunc("/healthz", check.HealthCheckHandler(cache, 29*time.Second)) // Тайминг можно настроить

	// Эндпоинт для readiness check
	mux.HandleFunc("/readyz", check.ReadinessCheckHandler(checker))

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
	h = handler.LoggingMiddleware(log, h)

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
		log.Info("http server listening", "addr", srv.Addr, "pid", os.Getpid())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	// Сигнал или ошибка сервера
	select {
	case <-rootCtx.Done():
		log.Info("shutdown signal received")
	case err := <-errCh:
		if err != nil {
			log.Error("server failed", "err", err)
			_ = cache.Close()
			return
		}
	}

	// Graceful shutdown

	graceTime := 30 * time.Second
	shCtx, cancel := context.WithTimeout(ctx, graceTime)
	defer cancel()

	if err := srv.Shutdown(shCtx); err != nil {
		log.Error("http shutdown error", "err", err)
	} else {
		log.Info("http server shutdown complete")
	}

	if err := cache.Close(); err != nil {
		log.Error("cache close error", "err", err)
	} else {
		log.Info("cache closed")
	}

	log.Info("graceful shutdown finished")

}
