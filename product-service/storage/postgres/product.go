package postgres

import (
	pb "github.com/product-service/genproto/product_service"

	"github.com/jmoiron/sqlx"
)

type productRepo struct {
	db *sqlx.DB
}

// NewProductRepo ...
func NewProductRepo(db *sqlx.DB) *productRepo {
	return &productRepo{db: db}
}

func (r *productRepo) CreateProduct(product *pb.Product) (*pb.Product, error) {
	productResp := pb.Product{}
	err := r.db.QueryRow(`
	insert into products	 (
		name,
		model,
		owner_id) values($1, $2, $3) returning 
		id, 
		name, 
		model, 
		owner_id`,
		product.Name,
		product.Model,
		product.OwnerId).
		Scan(
			&productResp.Id,
			&productResp.Name,
			&productResp.Model,
			&productResp.OwnerId,
		)
	if err != nil {
		return &pb.Product{}, err
	}
	return &productResp, nil
}

func (r *productRepo) GetProduct(ID int64) (*pb.Product, error) {
	product := pb.Product{}
	err := r.db.QueryRow(`select
	id,
	name,
	model,
	owner_id from products where id = $1`, ID).
		Scan(
			&product.Id,
			&product.Name,
			&product.Model,
			&product.OwnerId,
		)
	if err != nil {
		return &pb.Product{}, err
	}
	return &product, nil
}

func (r *productRepo) GetUserProducts(ownerID int64) (*pb.Products, error) {
	rows, err := r.db.Query(`select
	id,
	name,
	model,
	owner_id from products where owner_id = $1`, ownerID)

	if err != nil {
		return &pb.Products{}, err
	}
	defer rows.Close()
	products := &pb.Products{}

	for rows.Next() {
		product := &pb.Product{}
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Model,
			&product.OwnerId,
		)
		if err != nil {
			return &pb.Products{}, err
		}

		products.Products = append(products.Products, product)
		products.Count++
	}
	return products, nil
}
