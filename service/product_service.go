package service

import (
	"context"
	"test_project/model"
)

type (
	ProductService interface {
		CreateProduct(ctx context.Context, req *model.CreateProductRequest) error
		DeleteProduct(ctx context.Context, id string) error
		GetProduct(id string) (*model.DataProductResponse, error)
		ListProduct(req model.ListProductRequest) ([]model.DataProductResponse, int64, error)
		UpdateProduct(ctx context.Context, req *model.UpdateProductRequest) error
	}
)
