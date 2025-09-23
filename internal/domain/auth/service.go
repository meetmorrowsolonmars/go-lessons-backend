package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type Service struct {
	userService UserService
	jwtProvider JWTProvider
}

func NewService(
	userService UserService,
	jwtProvider JWTProvider,
) *Service {
	return &Service{
		userService: userService,
		jwtProvider: jwtProvider,
	}
}

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userService.GetByEmail(ctx, email)
	if errors.Is(err, model.ErrNotFound) {
		return "", fmt.Errorf("user not found: %w", model.ErrNotAuthorized)
	}
	if err != nil {
		return "", fmt.Errorf("get user by email: %w", err)
	}

	if !user.CheckPassword(password) {
		return "", fmt.Errorf("invalid password: %w", model.ErrNotAuthorized)
	}

	token, err := s.jwtProvider.Generate(user)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}
