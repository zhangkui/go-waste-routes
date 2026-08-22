package middleware

import (
	"net/http"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/transport/contextx"
	"go-waste-routes/internal/transport/http/response"
)

func RBAC(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := contextx.UserFromContext(r.Context())
			if !ok {
				response.Error(w, http.StatusForbidden, response.CodeForbidden, "forbidden", contextx.RequestID(r.Context()))
				return
			}
			if permission == "" || userHasPermission(user, permission) {
				next.ServeHTTP(w, r)
				return
			}
			response.Error(w, http.StatusForbidden, response.CodeForbidden, "forbidden", contextx.RequestID(r.Context()))
		})
	}
}

func userHasPermission(user domain.User, permission string) bool {
	for _, item := range user.Permissions {
		if item.Code == permission {
			return true
		}
	}
	return false
}
