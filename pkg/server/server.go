package server

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

const serviceName = "HTTP Server"

type server struct {
	httpServer *http.Server
	port       string
	host       string
	log        *zap.Logger
}

type Config struct {
	Host   string
	Port   string
	Router http.Handler
	Log    *zap.Logger
}

func New(config *Config) *server {
	return &server{
		port: config.Port,
		host: config.Host,
		log:  config.Log,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", config.Host, config.Port),
			Handler: config.Router,
		},
	}
}

func (s *server) Shutdown() error {
	return s.httpServer.Shutdown(context.Background())
}

func (s *server) ServiceName() string {
	return serviceName
}

func (s *server) Run() {
	s.log.Info("Listening HTTP", zap.String("host", s.host), zap.String("port", s.port))

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		s.log.Error("HTTP server ListenAndServe Error", zap.Error(err))
	}
}
