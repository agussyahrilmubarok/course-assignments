package client

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "example.com/pkg/proto/v1"
)

type UserGrpcClient struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
	logger *zap.Logger
}

func NewUserGrpcClient(grpcAddr string, logger *zap.Logger) (*UserGrpcClient, error) {
	conn, err := grpc.Dial(
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		logger.Error("failed to connect to grpc server", zap.Error(err))
		return nil, err
	}

	client := pb.NewUserServiceClient(conn)
	return &UserGrpcClient{conn: conn, client: client, logger: logger}, nil
}

func (c *UserGrpcClient) GetUserEmail(ctx context.Context, id string) (string, error) {
	resp, err := c.client.GetUserByID(ctx, &pb.GetUserRequest{Id: id})
	if err != nil {
		c.logger.Error("failed to get product by id", zap.String("id", id), zap.Error(err))
		return "", err
	}

	return resp.Email, nil
}

func (c *UserGrpcClient) Close() error {
	return c.conn.Close()
}
