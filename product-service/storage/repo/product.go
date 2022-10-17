package repo

import (
	pb "github.com/product-service/genproto/product_service"
)

//ProductStorageI ...
type ProductStorageI interface {
    CreateProduct(*pb.Product) (*pb.Product, error)
	GetProduct(ID int64) (*pb.Product, error)
	GetUserProducts(ownerID int64) (*pb.Products, error)
}