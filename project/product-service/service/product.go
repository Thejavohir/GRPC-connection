package service

import (
	"context"
	l "github.com/project/product-service/pkg/logger"
	pb "github.com/project/product-service/genproto/product_service"
	"github.com/project/product-service/storage"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductService struct {
	storage storage.IStorage
	logger  l.Logger
}

func NewProductService(db *sqlx.DB, log l.Logger) *ProductService {
	return &ProductService{
		storage: storage.NewStoragePg(db),
		logger:  log,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *pb.Product) (*pb.Product, error) {
	product, err := s.storage.Product().CreateProduct(req)
	if err != nil {
		s.logger.Error("error insert", l.Any("error insert product", err))
		return &pb.Product{}, status.Error(codes.Internal, "internal error")
	}
	return product, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	product, err := s.storage.Product().GetProduct(req.Id)
	if err != nil {
		s.logger.Error("error getting product", l.Any("error getting product", err))
		return &pb.Product{}, status.Error(codes.Internal, "internal error")
	}
	return product, nil
}

func (s *ProductService) GetUserProducts(ctx context.Context, req *pb.GetUserProductsRequest) (*pb.Products, error) {
	products, err := s.storage.Product().GetUserProducts(req.OwnerId)
	if err != nil {
		s.logger.Error("error getting user products", l.Any("error getting user products", err))
		return &pb.Products{}, status.Error(codes.Internal, "internal error")
	}
	return products, nil
}

func (s *ProductService) ListProducts(ctx context.Context, req *pb.LPreq) (*pb.LPresp, error) {
	products, err := s.storage.Product().ListProducts(req.Page, req.Limit)
	if err != nil {
		s.logger.Error("error getting user products", l.Any("error getting user products", err))
		return &pb.LPresp{}, status.Error(codes.Internal, "internal error")
	}
	return products, nil
}