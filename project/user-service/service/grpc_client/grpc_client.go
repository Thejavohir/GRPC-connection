package grpcClient

import (
	"fmt"

	"github.com/project/user-service/config"
	pbp "github.com/project/user-service/genproto/product_service"
	"google.golang.org/grpc"
)

// GrpcClientI ...
type GrpcClientI interface {
	Product() pbp.ProductServiceClient
}

// GrpcClient ...
type GrpcClient struct {
	cfg            config.Config
	productService pbp.ProductServiceClient
}

// New ...
func New(cfg config.Config) (*GrpcClient, error) {
	connProduct, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.ProductServiceHost, cfg.ProductServicePort),
		grpc.WithInsecure())

	if err != nil {
		return nil, fmt.Errorf("post service dial host: %s port: %d", cfg.ProductServiceHost, cfg.ProductServicePort)
	}
	return &GrpcClient{
		cfg:            cfg,
		productService: pbp.NewProductServiceClient(connProduct),
	}, nil
}

func (r *GrpcClient) Product() pbp.ProductServiceClient {
	return r.productService
}