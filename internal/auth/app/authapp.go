package authapp

import (
	"auth_records/pkg/server"
	"auth_records/pkg/shutdown"
	"auth_records/pkg/utils"

	"go.uber.org/zap"
)

func Run() {
	// Get ENV variables
	host := utils.GetEnv("HOST", "localhost")
	port := utils.GetEnv("PORT", "8080")

	// Initialize the logger
	logger, _ := zap.NewProduction()
	defer func() {
		err := logger.Sync()

		if err != nil {
			logger.Error("Error syncing logger", zap.Error(err))
		}
	}()

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
