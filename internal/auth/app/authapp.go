package authapp

import (
	"auth_records/internal/auth/adapter/api/v1"
	"auth_records/internal/auth/adapter/storage"
	"auth_records/internal/auth/infrastructure/pgprovider"
	"auth_records/internal/auth/infrastructure/recordsgrpc"
	"auth_records/internal/auth/infrastructure/router"
	"auth_records/internal/auth/usecase"
	"auth_records/pkg/server"
	"auth_records/pkg/shutdown"
	"auth_records/pkg/token"
	"auth_records/pkg/utils"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"

	"go.uber.org/zap"
)

func Run() {
	// Get ENV variables
	host := utils.GetEnv("HOST", "localhost")
	port := utils.GetEnv("PORT", "8080")

	dbUser := utils.GetEnv("DB_USER", "users_development")
	dbPassword := utils.GetEnv("DB_PASSWORD", "users_development")
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbDatabase := utils.GetEnv("DB_DATABASE", "auth_users_development")

	gRPCRecordsHost := utils.GetEnv("RECORDS_CLIENT_HOST", "records-service")
	gRPCRecordsPort := utils.GetEnv("RECORDS_CLIENT_PORT", "50051")

	secretKey := []byte(utils.GetEnv("JWT_SECRET_KEY", "aautreddf12w"))

	dbURL := utils.GetEnv("DATABASE_URL", fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbDatabase))

	// Initialize the logger
	logger, _ := zap.NewProduction()
	defer func() {
		err := logger.Sync()

		if err != nil {
			logger.Error("Error syncing logger", zap.Error(err))
		}
	}()

	// Initialize Token Service
	tokenService := token.NewJwtService(secretKey)

	// Initialize Pg Client
	pgClient := pgprovider.NewProvider(dbURL)
	err := pgClient.Connect()
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
	}

	recordsgRPCClient := recordsgrpc.New(logger, gRPCRecordsHost, gRPCRecordsPort)

	// Initialize Repository
	userRepo := storage.NewUserRepository(pgClient)

	// Initialize UseCase
	authUseCase := usecase.New(logger, userRepo, tokenService)

	// Initialize API Handler
	apiHandler := api.New(logger, authUseCase, recordsgRPCClient)

	// Initialize the Router
	appRouter := router.New(apiHandler)

	// Initialize the server
	s := server.New(&server.Config{
		Host:   host,
		Port:   port,
		Log:    logger,
		Router: appRouter,
	})

	// Initialize the shutdown watcher
	shutdownWatcher := shutdown.New(logger,
		s,
	)
	defer func() {
		err := shutdownWatcher.Shutdown()

		if err != nil {
			logger.Error("Error Shutdown service", zap.Error(err))
		}
	}()

	// Run the server
	go s.Run()

	// Wait for the signal
	shutdownWatcher.Wait()
}
