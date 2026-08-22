package contextx

import (
	"context"

	"go-waste-routes/internal/domain"
	appjwt "go-waste-routes/internal/platform/jwt"
)

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	claimsKey    contextKey = "claims"
	userKey      contextKey = "user"
)

func WithRequestID(ctx context.Context, requestID string) context.Context { return context.WithValue(ctx, requestIDKey, requestID) }
func RequestID(ctx context.Context) string { value, _ := ctx.Value(requestIDKey).(string); return value }
func WithClaims(ctx context.Context, claims *appjwt.Claims) context.Context { return context.WithValue(ctx, claimsKey, claims) }
func ClaimsFromContext(ctx context.Context) *appjwt.Claims { value, _ := ctx.Value(claimsKey).(*appjwt.Claims); return value }
func WithUser(ctx context.Context, user domain.User) context.Context { return context.WithValue(ctx, userKey, user) }
func UserFromContext(ctx context.Context) (domain.User, bool) { value, ok := ctx.Value(userKey).(domain.User); return value, ok }
