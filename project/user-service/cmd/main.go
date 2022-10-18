package main

import (
	"net"

	"github.com/project/user-service/config"
	pb "github.com/project/user-service/genproto/user_service"
	grpcClient "github.com/project/user-service/service/grpc_client"

	"github.com/project/user-service/pkg/db"
	"github.com/project/user-service/pkg/logger"
	"github.com/project/user-service/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "user-service")
	defer logger.Cleanup(log)

	log.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	conDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx failed to connect to database", logger.Error(err))
	}

	grpcClient, err := grpcClient.New(cfg)
	if err != nil {
		log.Fatal("error while gprc connection with client", logger.Error(err))
	}

	userService := service.NewUserService(conDB, log, grpcClient)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("error while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterUserServiceServer(s, userService)
	log.Info("main: server running", logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("error while listening %v", logger.Error(err))
	}
}
