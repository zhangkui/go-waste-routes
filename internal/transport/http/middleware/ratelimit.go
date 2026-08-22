package middleware

import (
	"net/http"
	"sync"
	"time"

	"go-waste-routes/internal/transport/contextx"
	"go-waste-routes/internal/transport/http/response"
)

func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
	var mu sync.Mutex
	var buckets = map[string][]time.Time{}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			now := time.Now()
			mu.Lock()
			bucket := buckets[ip]
			filtered := bucket[:0]
			for _, ts := range bucket {
				if now.Sub(ts) < window {
					filtered = append(filtered, ts)
				}
			}
			if len(filtered) >= limit {
				mu.Unlock()
				response.Error(w, http.StatusTooManyRequests, response.CodeBusiness, "rate limit exceeded", contextx.RequestID(r.Context()))
				return
			}
			buckets[ip] = append(filtered, now)
			mu.Unlock()
			next.ServeHTTP(w, r)
		})
	}
}
