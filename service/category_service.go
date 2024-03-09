package service

import "test_project/model"

type (
	CategoryService interface {
		CreateCategory(req *model.CreateCategoryRequest) error
		DeleteCategory(id string) error
		GetCategory(id string) (*model.DataCategoryResponse, error)
		ListCategory(req model.ListCategoryRequest) ([]model.DataCategoryResponse, int64, error)
		UpdateCategory(req *model.UpdateCategoryRequest) error
	}
)
