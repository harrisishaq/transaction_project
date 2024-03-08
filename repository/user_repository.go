package repository

import "test_project/entity"

type (
	UserRepository interface {
		Create(model entity.User) (string, error)
		Delete(model *entity.User) error
		Get(id string) (*entity.User, error)
		List(limit, offset int) ([]entity.User, int64, error)
		Update(model *entity.User) error
	}
)
