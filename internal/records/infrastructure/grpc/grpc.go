package grpc

import (
	"auth_records/pkg/grpcserver"
	"auth_records/pkg/middleware"

	pb "auth_records/pkg/records_grpc/v1"

	"go.uber.org/zap"
)

const serviceName = "gRPC Server"

type grpc struct {
	log    *zap.Logger
	port   string
	server *grpcserver.Server
}

func NewServer(log *zap.Logger, port string, getterHandlers pb.RecordsServiceServer, authMiddleware *middleware.AuthMiddleware) *grpc {
	s := grpcserver.New(log, port, authMiddleware)

	pb.RegisterRecordsServiceServer(s.Server(), getterHandlers)

	return &grpc{
		log:    log,
		port:   port,
		server: s,
	}
}

func (g *grpc) ServiceName() string {
	return serviceName
}

func (g *grpc) Run() {
	g.server.Run()
}

func (g *grpc) Shutdown() error {
	g.server.Shutdown()

	return nil
}
