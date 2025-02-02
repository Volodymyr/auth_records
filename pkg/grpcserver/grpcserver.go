package grpcserver

import (
	"fmt"
	"net"

	"auth_records/pkg/middleware"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	log    *zap.Logger
	port   string
	server *grpc.Server
}

func New(log *zap.Logger, port string, authMiddleware *middleware.AuthMiddleware) *Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(authMiddleware.JwtUnaryInterceptor),
	)

	return &Server{
		log:    log,
		port:   port,
		server: server,
	}
}

func (g *Server) Server() *grpc.Server {
	return g.server
}

func (g *Server) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", g.port))
	if err != nil {
		g.log.Error("failed to listen tcp", zap.Error(err))

		return
	}

	g.log.Info("server listening at", zap.String("address", lis.Addr().String()))
	if err := g.server.Serve(lis); err != nil {
		g.log.Error("failed to serve ", zap.Error(err))
	}
}

func (g *Server) Shutdown() {
	g.server.GracefulStop()
}
