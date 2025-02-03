package recordsgrpc

import (
	"auth_records/pkg/grpcclient"
	pb "auth_records/pkg/records_grpc/v1"
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const serviceName = "Records gRPC Client"

type recordsgrpc struct {
	getterClient pb.RecordsServiceClient
	conn         *grpc.ClientConn
}

func New(log *zap.Logger, host, port string) *recordsgrpc {
	recordsgRPCClient := grpcclient.New(serviceName, log, host, port)
	recordsgRPCClient.Connect()

	conn := recordsgRPCClient.Conn()
	getterClient := pb.NewRecordsServiceClient(conn)

	return &recordsgrpc{
		getterClient,
		conn,
	}
}

func (r *recordsgrpc) GetRandomRecords(token string) ([]*pb.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// add jwt token for metadata
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	req := &pb.GetRandomRecordsRequest{}

	resp, err := r.getterClient.GetRandomRecords(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Records, nil
}

func (r *recordsgrpc) ServiceName() string {
	return serviceName
}

func (r *recordsgrpc) Shutdown() error {
	return r.conn.Close()
}
