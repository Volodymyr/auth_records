package grpcclient

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	host        string
	port        string
	conn        *grpc.ClientConn
	log         *zap.Logger
	serviceName string
}

func New(serviceName string, log *zap.Logger, host, port string) *client {
	return &client{
		host:        host,
		port:        port,
		log:         log,
		serviceName: serviceName,
	}
}

func (c *client) Connect() {
	c.log.Info("Connecting gRPC", zap.String("name", c.serviceName), zap.String("host", c.host), zap.String("port", c.port))

	conn, err := grpc.NewClient(c.host+":"+c.port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.log.Error("Failed to connect to users gRPC client", zap.Error(err))
	}

	c.log.Info(conn.GetState().String())

	c.conn = conn
}

func (c *client) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *client) ServiceName() string {
	return c.serviceName
}

func (c *client) Shutdown() error {
	return c.conn.Close()
}
