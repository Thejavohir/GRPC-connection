package service

import (
	"context"

	pb "github.com/project/user-service/genproto/user_service"
	l "github.com/project/user-service/pkg/logger"
	grpcClient "github.com/project/user-service/service/grpc_client"
	"github.com/project/user-service/storage"

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
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	user, err := s.storage.User().UpdateUser(req)
	if err != nil {
		s.logger.Error("error getting user products", l.Any("error getting user products", err))
		return &pb.User{}, status.Error(codes.Internal, "internal error")
	}
	return &pb.User{
		Id: user.Id,
	}, nil
}

func (s *UserService) CheckField(ctx context.Context, req *pb.CheckFieldRequest) (*pb.CheckFieldResponse, error) {
	boolean, err := s.storage.User().CheckField(req.Field, req.Value)
	if err != nil {
		s.logger.Error("error checkfield user", l.Any("error checking field", err))
		return &pb.CheckFieldResponse{}, status.Error(codes.Internal, "internal error")
	}
	return &pb.CheckFieldResponse{
		Exists: boolean.Exists,
	}, nil
}