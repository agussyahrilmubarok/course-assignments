package client

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "example.com/pkg/proto/v1"
)

type ProductGrpcClient struct {
	conn   *grpc.ClientConn
	client pb.ProductServiceClient
	logger *zap.Logger
}

func NewProductGrpcClient(grpcAddr string, logger *zap.Logger) (*ProductGrpcClient, error) {
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

	client := pb.NewProductServiceClient(conn)
	return &ProductGrpcClient{conn: conn, client: client, logger: logger}, nil
}

func (c *ProductGrpcClient) GetProductInfo(ctx context.Context, id string) (string, float64, int, error) {
	resp, err := c.client.GetProductByID(ctx, &pb.GetProductRequest{Id: id})
	if err != nil {
		c.logger.Error("failed to get product by id", zap.String("id", id), zap.Error(err))
		return "", 0, 0, err
	}

	return resp.Name, resp.Price, int(resp.Stock), nil
}

func (c *ProductGrpcClient) Close() error {
	return c.conn.Close()
}
