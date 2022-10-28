package repo

import pb "github.com/project/user-service/genproto/user_service"

type UserStorageI interface {
	Create(*pb.User) (*pb.User, error)
	GetUser(ID int64) (*pb.User, error)
	UpdateUser(*pb.User) (*pb.User, error)
	CheckField(field, value string) (*pb.CheckFieldResponse, error)  
}