package repository

import "test_project/entity"

type (
	ProductRepository interface {
		Create(model entity.Product) (string, error)
		Delete(model *entity.Product) error
		Get(id string) (*entity.Product, error)
		List(limit, offset int, filters map[string]interface{}) ([]entity.Product, int64, error)
		Update(model *entity.Product) error
		// Log
		CreateLog(model entity.ProductLog) error
	}
)
