package repository

import "test_project/entity"

type (
	CategoryRepository interface {
		Create(model entity.Category) (string, error)
		Delete(model *entity.Category) error
		Get(id string) (*entity.Category, error)
		List(limit, offset int, filters map[string]interface{}) ([]entity.Category, int64, error)
		Update(model *entity.Category) error
		// Log
		CreateLog(model entity.CategoryLog) error
	}
)
