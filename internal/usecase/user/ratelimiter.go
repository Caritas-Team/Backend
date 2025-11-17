package user

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Caritas-Team/reviewer/internal/config"
	"github.com/Caritas-Team/reviewer/internal/memcached"
)

var ErrRateLimitExceeded = errors.New("too many requests")

type RateLimiter struct {
	cache    memcached.CacheInterface
	enabled  bool
	window   time.Duration
	requests int
}

func NewRateLimiter(cache memcached.CacheInterface, cfg config.Config) *RateLimiter {
	return &RateLimiter{
		cache:    cache,
		enabled:  cfg.RateLimiter.Enabled,
		window:   time.Duration(cfg.RateLimiter.WindowSize) * time.Second,
		requests: cfg.RateLimiter.RequestsPerWindow,
	}
}

func (rl *RateLimiter) AllowRequest(ctx context.Context, userID string) error {
	if !rl.enabled {
		return nil
	}

	key := "rate_limit:" + userID

	newValue, err := rl.cache.Increment(ctx, key, 1)
	if err != nil {
		return nil
	}

	// Устанавливаем TTL при первом запросе
	if newValue == 1 {
		err := rl.cache.Set(ctx, key, []byte("1"), rl.window)
		if err != nil {
			return err
		}
	}

	if newValue > uint64(rl.requests) {
		_, err := rl.cache.Decrement(ctx, key, 1)
		if err != nil {
			slog.Warn("Failed to decrement rate limit counter", "error", err)
		}
		return ErrRateLimitExceeded
	}

	return nil
}
