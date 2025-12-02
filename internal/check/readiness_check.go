package check

import (
	"log"
	"net/http"

	"github.com/Caritas-Team/reviewer/internal/handler"
	"github.com/Caritas-Team/reviewer/internal/memcached"
	"github.com/Caritas-Team/reviewer/internal/metrics"
	"github.com/Caritas-Team/reviewer/internal/logger"
)

log := logger.NewLogger(cfg)

// ReadinessChecker проверяет состояние приложения
type ReadinessChecker struct {
	cache       *memcached.Cache
	rateLimiter *handler.RateLimiterMiddleware
}

// Конструктор ReadinessChecker
func NewReadinessChecker(cache *memcached.Cache, rateLimiter *handler.RateLimiterMiddleware) *ReadinessChecker {
	return &ReadinessChecker{
		cache:       cache,
		rateLimiter: rateLimiter,
	}
}

// Проверка готовности
func (rc *ReadinessChecker) IsReady() bool {

	// Проверка memcached
	if err := rc.cache.IsEnabled(); err != nil {
		log.Error("Ошибка проверки активности кэша:", err)
		return false
	}

	if err := rc.cache.Ping(); err != nil {
		return false
	}

	// Проверка CORS (конфигурация передаётся в функцию)
	if err := handler.CheckCORS(handler.CORSConfig{}); err != nil {
		return false
	}

	// Проверка rate limiting
	if err := rc.rateLimiter.IsOperational(); err != nil {
		return false
	}

	// Проверка метрик
	if err := metrics.CheckMetrics(); err != nil {
		return false
	}

	// Заглушки для будущих методов
	if !uploadStub() {
		return false
	}

	if !processingStub() {
		return false
	}

	return true
}

// Заглушка для метода upload
func uploadStub() bool {

	return true
}

// Заглушка для метода processing
func processingStub() bool {

	return true
}

// Обработчик readiness check
func ReadinessCheckHandler(checker *ReadinessChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if checker.IsReady() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("READY"))
		} else {
			http.Error(w, "NOT READY", http.StatusServiceUnavailable)
		}
	}
}
