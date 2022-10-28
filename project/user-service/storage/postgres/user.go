package postgres

import (
	"fmt"

	pb "github.com/project/user-service/genproto/user_service"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

// NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *pb.User) (*pb.User, error) {
	userResp := pb.User{}
	err := r.db.QueryRow(`
	insert into users(
		name,
		last_name,
		username,
		email,
		password) 
		values ($1, $2, $3, $4, $5) 
		returning id, name, last_name, username, email, password`,
		user.Name,
		user.Lastname,
		user.Username,
		user.Email,
		user.Password).Scan(
		&userResp.Id,
		&userResp.Name,
		&userResp.Lastname,
		&userResp.Username,
		&userResp.Email,
		&userResp.Password,
	)
	if err != nil {
		return &pb.User{}, err
	}
	return &userResp, nil
}

func (r *userRepo) GetUser(ID int64) (*pb.User, error) {
	user := pb.User{}
	err := r.db.QueryRow(`select 
	id, 
	name, 
	last_name from users where id = $1`, ID).Scan(
		&user.Id,
		&user.Name,
		&user.Lastname,
	)
	if err != nil {
		return &pb.User{}, err
	}
	return &user, nil
}

func (r *userRepo) UpdateUser(user *pb.User) (*pb.User, error) {
	userUp := pb.User{}
	err := r.db.QueryRow(`insert into users(
	id,
	name,
	last_name) values($1, $2, $3) returning 
	id, 
	name, 
	last_name`,
		user.Id,
		user.Name,
		user.Lastname).Scan(
		&userUp.Id,
	)
	if err != nil {
		return &pb.User{}, err
	}
	return &userUp, nil
}

func (r *userRepo) CheckField(field, value string) (*pb.CheckFieldResponse, error) {
	query := fmt.Sprintf("select count(1) from users where %s = $1", field)
	var exists int
	err := r.db.QueryRow(query, value).Scan(&exists)
	if err != nil {
		return &pb.CheckFieldResponse{}, err
	}
	if exists == 0 {
		return &pb.CheckFieldResponse{Exists: false}, nil
	}
	return &pb.CheckFieldResponse{Exists: true}, nil
}