package recordapp

import (
	"auth_records/internal/records/adapters/grpcgetterhandlers"
	"auth_records/internal/records/adapters/storage"
	"auth_records/internal/records/infrastructure/grpc"
	"auth_records/internal/records/infrastructure/pgprovider"
	"auth_records/internal/records/usecase"
	"auth_records/pkg/middleware"
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
	grpcPort := utils.GetEnv("GRPC_PORT", "50051")

	dbUser := utils.GetEnv("DB_USER", "records_development")
	dbPassword := utils.GetEnv("DB_PASSWORD", "records_development")
	dbHost := utils.GetEnv("DB_HOST", "localhost")
	dbPort := utils.GetEnv("DB_PORT", "5432")
	dbDatabase := utils.GetEnv("DB_DATABASE", "records_development")

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

	// Initialize repositories
	recordsRepo := storage.NewRecordsRepository(pgClient)

	// Initialize use case
	recordsUseCase := usecase.New(logger, recordsRepo)

	// Initialize grpc handler
	grpcServerHandler := grpcgetterhandlers.New(logger, recordsUseCase)

	jwtService := token.NewJwtService(secretKey)
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	//Initialize grpc server
	grpcServer := grpc.NewServer(logger, grpcPort, grpcServerHandler, authMiddleware)

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
