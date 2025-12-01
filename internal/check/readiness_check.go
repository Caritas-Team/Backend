package check

import (
	"net/http"

	"github.com/Caritas-Team/reviewer/internal/handler"
	"github.com/Caritas-Team/reviewer/internal/logger"
	"github.com/Caritas-Team/reviewer/internal/memcached"
	"github.com/Caritas-Team/reviewer/internal/metrics"
)

// ReadinessChecker проверяет состояние приложения
type ReadinessChecker struct {
	cache       *memcached.Cache
	rateLimiter *handler.RateLimiterMiddleware
	logger      *logger.Logger
}

// Конструктор ReadinessChecker
func NewReadinessChecker(cache *memcached.Cache, rateLimiter *handler.RateLimiterMiddleware, logger *logger.Logger) *ReadinessChecker {
	return &ReadinessChecker{
		cache:       cache,
		rateLimiter: rateLimiter,
		logger:      logger,
	}
}

// Проверка готовности
func (rc *ReadinessChecker) IsReady() bool {
	// Проверка memcached
	if !rc.cache.IsEnabled() || rc.cache.Ping() != nil {
		return false
	}

	// Проверка CORS (конфигурация передаётся в функцию)
	if !handler.CheckCORS(handler.CORSConfig{}) {
		return false
	}

	// Проверка rate limiting
	if !rc.rateLimiter.IsOperational() {
		return false
	}

	// Проверка метрик
	if !metrics.CheckMetrics() {
		return false
	}

	// Проверка логгера
	if err := rc.logger.TestLog(); err != nil {
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
