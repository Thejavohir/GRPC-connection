package postgres

import (
	pb "github.com/user-service/genproto/user_service"

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
		last_name) 
		values ($1, $2) 
		returning 
		id, 
		name, 
		last_name`, 
		user.Id, 
		user.Name, 
		user.Lastname).Scan(
			&userResp.Id,
			&userResp.Name,
			&userResp.Lastname,
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