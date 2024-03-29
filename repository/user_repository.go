package repository

import "test_project/entity"

type (
	UserRepository interface {
		Create(model entity.User) (string, error)
		Delete(model *entity.User) error
		Get(id string) (*entity.User, error)
		GetByEmail(email string) (*entity.User, error)
		List(limit, offset int, filters map[string]interface{}) ([]entity.User, int64, error)
		Update(model *entity.User) error
		// Log
		CreateLog(model entity.UserLog) error
	}
)
