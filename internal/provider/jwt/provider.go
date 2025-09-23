package jwt

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type Claims struct {
	model.AuthClaims
	jwt.RegisteredClaims
}

type Provider struct {
	secretKey           []byte
	issuer              string
	accessTokenDuration time.Duration

	now func() time.Time
}

func NewProvider(
	secretKey []byte,
	issuer string,
	accessTokenDuration time.Duration,
) *Provider {
	return &Provider{
		secretKey:           secretKey,
		issuer:              issuer,
		accessTokenDuration: accessTokenDuration,
		now:                 func() time.Time { return time.Now().UTC() },
	}
}

func (p *Provider) Generate(user model.User) (string, error) {
	now := p.now()
	subject := strconv.FormatInt(user.ID, 10)

	claims := Claims{
		AuthClaims: model.AuthClaims{
			UserID: user.ID,
			Email:  user.Email,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    p.issuer,
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(now.Add(p.accessTokenDuration)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	// TODO: Use RS256 instead HS256.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(p.secretKey)
}

func (p *Provider) Validate(tokenString string) (model.AuthClaims, error) {
	var claims Claims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return p.secretKey, nil
		}

		return nil, model.ErrNotAuthorized
	})
	if err != nil {
		return model.AuthClaims{}, err
	}

	if !token.Valid {
		return model.AuthClaims{}, model.ErrNotAuthorized
	}

	return claims.AuthClaims, nil
}
