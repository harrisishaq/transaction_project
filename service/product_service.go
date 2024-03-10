package service

import "test_project/model"

type (
	ProductService interface {
		CreateProduct(req *model.CreateProductRequest) error
		DeleteProduct(id string) error
		GetProduct(id string) (*model.DataProductResponse, error)
		ListProduct(req model.ListProductRequest) ([]model.DataProductResponse, int64, error)
		UpdateProduct(req *model.UpdateProductRequest) error
	}
)
