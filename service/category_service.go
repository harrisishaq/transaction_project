package service

import (
	"context"
	"test_project/model"
)

type (
	CategoryService interface {
		CreateCategory(ctx context.Context, req *model.CreateCategoryRequest) error
		DeleteCategory(ctx context.Context, id string) error
		GetCategory(id string) (*model.DataCategoryResponse, error)
		ListCategory(req model.ListCategoryRequest) ([]model.DataCategoryResponse, int64, error)
		UpdateCategory(ctx context.Context, req *model.UpdateCategoryRequest) error
	}
)
