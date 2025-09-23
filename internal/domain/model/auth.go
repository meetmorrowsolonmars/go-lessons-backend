package model

import (
	"context"
)

type AuthClaims struct {
	UserID int64
	Email  string
}

type ctxKey string

const authClaimsContextKey ctxKey = "auth_claims"

func WithAuthClaims(ctx context.Context, claims AuthClaims) context.Context {
	return context.WithValue(ctx, authClaimsContextKey, claims)
}

func AuthClaimsValue(ctx context.Context) (AuthClaims, bool) {
	claims, ok := ctx.Value(authClaimsContextKey).(AuthClaims)
	return claims, ok
}
