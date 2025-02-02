package recordapp

import (
	"auth_records/internal/records/infrastructure/grpc"
	"auth_records/internal/records/infrastructure/pgprovider"
	"auth_records/pkg/middleware"
	"auth_records/pkg/server"
	"auth_records/pkg/shutdown"
	"auth_records/pkg/token"
	"auth_records/pkg/utils"
	"fmt"

	"go.uber.org/zap"
)

func Run() {
	// Get ENV variables
	host := utils.GetEnv("HOST", "localhost")
	port := utils.GetEnv("PORT", "8080")

	dbUser := utils.GetEnv("DB_USER", "records_development")
	dbPassword := utils.GetEnv("DB_PASSWORD", "records_development")
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbDatabase := utils.GetEnv("DB_DATABASE", "records_development")

	//grpcPort := utils.GetEnv("USERS_SERVER_PORT", "50051")

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

	// Initialize Pg Client
	pgClient := pgprovider.NewProvider(dbURL)
	err := pgClient.Connect()
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
	}

	jwtService := token.NewJwtService(secretKey)
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	grpcServer := grpc.NewServer(logger, port, grpcServerHandler, authMiddleware)

	// Initialize the server
	s := server.New(&server.Config{
		Host: host,
		Port: port,
		Log:  logger,
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

	// Run grpc server
	go grpcServer.Run()

	// Run the server
	go s.Run()

	// Wait for the signal
	shutdownWatcher.Wait()
}
