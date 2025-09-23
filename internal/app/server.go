package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	accountv1 "github.com/meetmorrowsolonmars/education-pet-project/internal/api/v1/account"
	authv1 "github.com/meetmorrowsolonmars/education-pet-project/internal/api/v1/auth"
	userv1 "github.com/meetmorrowsolonmars/education-pet-project/internal/api/v1/user"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/auth"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/user"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/provider/jwt"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/provider/memory"
)

func RunServer() error {
	// Configure logger.
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}),
	)

	// Read app configuration.
	config, err := ReadConfig()
	if err != nil {
		logger.Error("Read config", slog.String("error", err.Error()))
		return fmt.Errorf("read config: %w", err)
	}

	// Configure providers.
	jwtProvider := jwt.NewProvider(config.JWT.SecretKey, config.JWT.Issuer, config.JWT.AccessTokenDuration)

	// Configure stores.
	userStore := memory.NewUserStore()
	accountStore := memory.NewAccountStore()

	// Configure services.
	userService := user.NewService(userStore, accountStore)
	authService := auth.NewService(userService, jwtProvider)

	// Configure controllers.
	authHandler := authv1.NewHandler(authService, logger)
	accountHandler := accountv1.NewHandler(accountStore, jwtProvider, logger)
	userHandler := userv1.NewHandler(userService, jwtProvider, logger)

	// Configure a HTTP server.
	mux := http.NewServeMux()

	authHandler.Register(mux)
	accountHandler.Register(mux)
	userHandler.Register(mux)

	server := &http.Server{
		Addr:    config.Server.Address,
		Handler: mux,
	}

	// Graceful shutdown.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Start the HTTP server.
	go func() {
		logger.Info("Start HTTP server", slog.String("address", server.Addr))

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server error", slog.String("error", err.Error()))
			cancel()
		}
	}()

	select {
	case <-ctx.Done():
	}

	logger.Info("Stop HTTP server", slog.String("address", server.Addr))

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = server.Shutdown(ctx)

	logger.Info("Server stopped", slog.String("address", server.Addr))

	return nil
}
