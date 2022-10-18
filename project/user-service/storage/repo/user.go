package repo

import pb "github.com/project/user-service/genproto/user_service"

type UserStorageI interface {
	Create(*pb.User) (*pb.User, error)
	GetUser(ID int64) (*pb.User, error)
	Update(user *pb.User) (*pb.User, error)
	// Delete(userID *pb.User) error
}