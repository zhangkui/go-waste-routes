package middleware

import (
	"net/http"

	"go-waste-routes/internal/transport/contextx"
	"go-waste-routes/internal/transport/http/response"
)

func Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					response.Error(w, http.StatusInternalServerError, response.CodeSystemError, "internal server error", contextx.RequestID(r.Context()))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
