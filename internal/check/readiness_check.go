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
	log         *logger.Logger
}

// Конструктор ReadinessChecker
func NewReadinessChecker(cache *memcached.Cache, rateLimiter *handler.RateLimiterMiddleware, log *logger.Logger) *ReadinessChecker {
	return &ReadinessChecker{
		cache:       cache,
		rateLimiter: rateLimiter,
		log:         log,
	}
}

// Проверка готовности
func (rc *ReadinessChecker) IsReady() bool {

	// Проверка memcached
	if err := rc.cache.IsEnabled(); err != nil {
		rc.log.Error("Ошибка проверки активности кэша:", err)
		return false
	}

	if err := rc.cache.Ping(); err != nil {
		rc.log.Error("Ошибка проверки доступности кэша:", err)
		return false
	}

	// Проверка CORS (конфигурация передаётся в функцию)
	if err := handler.CheckCORS(handler.CORSConfig{}); err != nil {
		rc.log.Error("Ошибка проверки CORS:", err)
		return false
	}

	// Проверка rate limiting
	if err := rc.rateLimiter.IsOperational(); err != nil {
		rc.log.Error("Ошибка проверки rate limiting:", err)
		return false
	}

	// Проверка метрик
	if err := metrics.CheckMetrics(); err != nil {
		rc.log.Error("Ошибка проверки метрик:", err)
		return false
	}

	// Заглушки для будущих методов
	if !uploadStub() {
		rc.log.Error("Ошибка проверки uploadStub")
		return false
	}

	if !processingStub() {
		rc.log.Error("Ошибка проверки processingStub")
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
			JSONResponse(w, http.StatusOK, "READY")
		} else {
			JSONResponse(w, http.StatusServiceUnavailable, "NOT READY")
		}
	}
}
