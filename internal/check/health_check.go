package check

import (
	"net/http"
	"time"

	"github.com/Caritas-Team/reviewer/internal/memcached"
)

// HealthCheckHandler — обработка health-check запроса с проверкой memcached и таймингом
func HealthCheckHandler(cache *memcached.Cache, maxResponseTime time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			elapsed := time.Since(start)
			if elapsed > maxResponseTime {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("Health check time out"))
				return
			}

			// Проверка доступности memcached
			if err := cache.IsEnabled(); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte("Memcached неактивен"))
				return
			}

			if err := cache.Ping(); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte("Memcached недоступен"))
				return
			}

			// Самотест
			// Выполняем простейший запрос к самому себе
			resp, err := http.Get(r.URL.Scheme + "://" + r.Host + "/ping")
			if err != nil || resp.StatusCode != http.StatusOK {
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte("Самотест не удался"))
				return
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		}()
	}
}
