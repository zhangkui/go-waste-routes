package middleware

import (
	"context"
	"net/http"
	"strings"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/platform/jwt"
	"go-waste-routes/internal/transport/contextx"
	"go-waste-routes/internal/transport/http/response"
)

type AuthProvider interface {
	CheckAccess(string) (*jwt.Claims, error)
	Me(context.Context, int64) (domain.User, error)
}

func Auth(provider AuthProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if !strings.HasPrefix(strings.ToLower(authorization), "bearer ") {
				response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "missing token", contextx.RequestID(r.Context()))
				return
			}
			claims, err := provider.CheckAccess(strings.TrimSpace(authorization[7:]))
			if err != nil {
				response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "invalid token", contextx.RequestID(r.Context()))
				return
			}
			user, err := provider.Me(r.Context(), claims.UserID)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, response.CodeUnauthorized, "invalid user", contextx.RequestID(r.Context()))
				return
			}
			ctx := contextx.WithClaims(r.Context(), claims)
			ctx = contextx.WithUser(ctx, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
