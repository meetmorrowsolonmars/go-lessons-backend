package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/examples/middleware/httpmiddleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/api/middleware"
	accountv1 "github.com/meetmorrowsolonmars/education-pet-project/internal/api/v1/account"
	authv1 "github.com/meetmorrowsolonmars/education-pet-project/internal/api/v1/auth"
	operationv1 "github.com/meetmorrowsolonmars/education-pet-project/internal/api/v1/operation"
	userv1 "github.com/meetmorrowsolonmars/education-pet-project/internal/api/v1/user"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/auth"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/operation"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/user"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/metric"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/provider/jwt"
	"github.com/meetmorrowsolonmars/education-pet-project/internal/provider/postgres"
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
	const envVarConfigPath = "CONFIG_PATH"
	configPath := os.Getenv(envVarConfigPath)
	if configPath == "" {
		logger.Error("Config path must be set")
		return fmt.Errorf("config path %s must be set", envVarConfigPath)
	}

	config, err := ReadConfig(configPath)
	if err != nil {
		logger.Error("Read config", slog.String("error", err.Error()))
		return fmt.Errorf("read config: %w", err)
	}

	// Configure database.
	dbConfig, err := pgxpool.ParseConfig(config.Postgres.ConnectionString)
	if err != nil {
		logger.Error("Parse db config", slog.String("error", err.Error()))
		return fmt.Errorf("parse db config: %w", err)
	}

	// TODO: Use normal context.
	pool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		logger.Error("Create db connection", slog.String("error", err.Error()))
		return fmt.Errorf("create db connection: %w", err)
	}

	// Configure traces.

	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)

	otlpTraceExporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint("localhost:4317"),
		otlptracehttp.WithTimeout(5*time.Second),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		logger.Error("Create OTLP trace exporter", slog.String("error", err.Error()))
		return fmt.Errorf("create OTLP trace exporter: %w", err)
	}

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(otlpTraceExporter,
		sdktrace.WithBatchTimeout(5*time.Second),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(batchSpanProcessor),
	)

	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagator)

	// TODO: use for graceful shutdown.
	// tracerProvider.Shutdown(context.Background())

	otel.GetTracerProvider()

	tracer := otel.Tracer("")
	_, span := tracer.Start(context.Background(), "test")
	span.AddEvent("hello world")
	span.End()

	// Configure metrics.

	prometheus.NewRegistry()

	registry := prometheus.NewRegistry()

	registry.MustRegister(
		collectors.NewGoCollector(
			collectors.WithGoCollectorRuntimeMetrics(
				collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/sched/latencies:seconds")},
			),
		),
	)

	metric.MustRegister(registry)

	promMiddleware := httpmiddleware.New(registry, prometheus.DefBuckets)

	// Configure providers.
	jwtProvider := jwt.NewProvider([]byte(config.JWT.SecretKey), config.JWT.Issuer, config.JWT.AccessTokenDuration)

	// Configure stores.
	userStore := postgres.NewUserStore(pool)
	accountStore := postgres.NewAccountStore(pool)
	operationStore := postgres.NewOperationStore(pool)
	categoryStore := postgres.NewCategoryStore(pool)

	// Configure services.
	userService := user.NewService(userStore)
	authService := auth.NewService(userService, jwtProvider)
	operationService := operation.NewService(operationStore, userStore, accountStore, categoryStore)

	// Configure controllers.
	authHandler := authv1.NewHandler(authService, logger)
	accountHandler := accountv1.NewHandler(accountStore, logger)
	userHandler := userv1.NewHandler(userService, logger)
	operationHandler := operationv1.NewHandler(operationService, logger)

	// Configure a HTTP server.
	authMiddleware := middleware.NewAuthMiddleware(jwtProvider)

	mux := http.NewServeMux()

	authHandler.Register(mux, promMiddleware, authMiddleware)
	accountHandler.Register(mux, promMiddleware, authMiddleware)
	userHandler.Register(mux, promMiddleware, authMiddleware)
	operationHandler.Register(mux, promMiddleware, authMiddleware)

	apiService := &http.Server{
		Addr:    config.Server.Address,
		Handler: mux,
	}

	mux = http.NewServeMux()

	mux.Handle("GET /metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	mux.HandleFunc("GET /debug/pprof/", pprof.Index)
	mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)

	debugServer := &http.Server{
		Addr:    config.DebugServer.Address,
		Handler: mux,
	}

	// Graceful shutdown.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Start the HTTP server.
	go func() {
		logger.Info("Start HTTP server", slog.String("address", apiService.Addr))

		if err := apiService.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server error", slog.String("error", err.Error()))
			cancel()
		}
	}()

	go func() {
		logger.Info("Start debug server", slog.String("address", debugServer.Addr))

		if err := debugServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server error", slog.String("error", err.Error()))
			cancel()
		}
	}()

	select {
	case <-ctx.Done():
	}

	logger.Info("Stop HTTP server", slog.String("address", apiService.Addr))

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = apiService.Shutdown(ctx)
	_ = debugServer.Shutdown(ctx)

	logger.Info("Server stopped", slog.String("address", apiService.Addr))

	return nil
}
