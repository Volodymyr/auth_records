package shutdown

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type GracefulShutdownService interface {
	Shutdown() error
	ServiceName() string
}

type watcher struct {
	log      *zap.Logger
	services []GracefulShutdownService
}

func New(log *zap.Logger, services ...GracefulShutdownService) *watcher {
	return &watcher{
		log,
		services,
	}
}

func (w *watcher) Shutdown() error {
	for _, service := range w.services {
		err := service.Shutdown()
		if err != nil {
			w.log.Error(fmt.Sprintf("Error shutting down %s: %s", service.ServiceName(), err.Error()))

			continue
		}

		w.log.Info(fmt.Sprintf("Shutdown %s. Bye bye.", service.ServiceName()))
	}

	return nil
}

func (w *watcher) WaitSignalAndShutdown() error {
	w.Wait()

	return w.Shutdown()
}

func (w *watcher) Wait() {
	sigint := make(chan os.Signal, 4)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-sigint
}
