package service

import (
	"context"

	pbp "github.com/user-service/genproto/product_service"
	pb "github.com/user-service/genproto/user_service"
	l "github.com/user-service/pkg/logger"
	grpcClient "github.com/user-service/service/grpc_client"
	"github.com/user-service/storage"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	storage storage.IStorage
	logger  l.Logger
	client  grpcClient.GrpcClientI
}

func NewUserService(db *sqlx.DB, log l.Logger, client grpcClient.GrpcClientI) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}

func (s *UserService) Create(ctx context.Context, req *pb.User) (*pb.User, error) {
	user, err := s.storage.User().Create(req)
	if err != nil {
		s.logger.Error("error insert", l.Any("error insert user", err))
		return &pb.User{}, status.Error(codes.Internal, "internal error")
	}
	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	user, err := s.storage.User().GetUser(req.Id)
	if err != nil {
		s.logger.Error("error getting user", l.Any("error getting user", err))
		return &pb.User{}, status.Error(codes.Internal, "internal error")
	}

	products, err := s.client.Product().GetUserProducts(ctx, &pbp.GetUserProductsRequest{
		OwnerId: req.Id,
	})
	if err != nil {
		s.logger.Error("error getting user products", l.Any("error getting user products", err))
		return &pb.User{}, status.Error(codes.Internal, "internal error")
	}
	for _, p := range products.Products {
		user.Products = append(user.Products, &pb.Product{
			Id:      p.Id,
			Name:    p.Name,
			Model:   p.Model,
			OwnerId: p.OwnerId,
		})
	}
	user.ProductsCount = products.Count

	return user, nil
}

// func (s *UserService) Update(ctx context.Context, req *pb.GetUserProductsRequest) (*pb.User, error) {
// 	products, err := s.storage.User().User(req.OwnerId)
// 	if err != nil {
// 		s.logger.Error("error getting user products", l.Any("error getting user products", err))
// 		return &pb.Products{}, status.Error(codes.Internal, "internal error")
// 	}
// 	return products, nil
// }
