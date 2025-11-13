package service

import (
	"context"

	pb "example.com/pkg/proto/v1"
	"example.com/productfc/internal/store"
	"go.uber.org/zap"
)

type productGrpcService struct {
	pb.UnimplementedProductServiceServer
	productStore store.IProductStore
	log          *zap.Logger
}

func NewProductGrpcService(productStore store.IProductStore, log *zap.Logger) *productGrpcService {
	return &productGrpcService{
		productStore: productStore,
		log:          log,
	}
}

func (s *productGrpcService) GetProductByID(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	product, err := s.productStore.FindByID(ctx, req.Id)
	if err != nil || product == nil {
		s.log.Error("failed to find product by id", zap.String("product_id", req.Id), zap.Error(err))
		return nil, err
	}

	return &pb.ProductResponse{
		Id:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: int32(product.Stock),
	}, nil
}
