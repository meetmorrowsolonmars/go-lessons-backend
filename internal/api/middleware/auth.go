package middleware

import (
	"net/http"
	"strings"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/api"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

//go:generate minimock -i JWTProvider

type JWTProvider interface {
	Validate(token string) (model.AuthClaims, error)
}

type AuthMiddleware interface {
	// WrapHandler wraps the given HTTP handler for authorization.
	WrapHandler(handler http.Handler) http.HandlerFunc
}

type authMiddleware struct {
	jwtProvider JWTProvider
}

func NewAuthMiddleware(jwtProvider JWTProvider) AuthMiddleware {
	return &authMiddleware{
		jwtProvider: jwtProvider,
	}
}

func (m *authMiddleware) WrapHandler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			api.EncodeErrorf(w, http.StatusUnauthorized, "Authorization header is required")

				return
			}

			parts := strings.SplitN(authorization, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				api.EncodeErrorf(w, http.StatusUnauthorized, "Authorization header format is invalid")

				return
			}

			token := parts[1]

		claims, err := m.jwtProvider.Validate(token)
		if err != nil {
			api.EncodeErrorf(w, http.StatusUnauthorized, "Validate token: %s", err)

				return
			}

		// TODO: Check that user exists.

		ctx := model.WithAuthClaims(r.Context(), claims)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)
	}
}
