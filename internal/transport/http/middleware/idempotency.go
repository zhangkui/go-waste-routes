package middleware

import (
	"net/http"
	"sync"
	"time"

	"go-waste-routes/internal/transport/contextx"
	"go-waste-routes/internal/transport/http/response"
)

func Idempotency(ttl time.Duration) func(http.Handler) http.Handler {
	var mu sync.Mutex
	var keys = map[string]time.Time{}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}
			now := time.Now()
			mu.Lock()
			for k, expiresAt := range keys {
				if now.After(expiresAt) {
					delete(keys, k)
				}
			}
			if expiresAt, ok := keys[key]; ok && now.Before(expiresAt) {
				mu.Unlock()
				response.Error(w, http.StatusConflict, response.CodeBusiness, "duplicate request", contextx.RequestID(r.Context()))
				return
			}
			keys[key] = now.Add(ttl)
			mu.Unlock()
			next.ServeHTTP(w, r)
		})
	}
}
