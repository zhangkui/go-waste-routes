package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go-waste-routes/internal/transport/contextx"
)

func Logger(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			next.ServeHTTP(w, r)
			log.WithFields(logrus.Fields{
				"request_id": contextx.RequestID(r.Context()),
				"method":     r.Method,
				"path":       r.URL.Path,
				"elapsed_ms": time.Since(started).Milliseconds(),
			}).Info("request")
		})
	}
}
