package repository

import "test_project/entity"

type (
	CustomerRepository interface {
		Create(model entity.Customer) (string, error)
		Delete(model *entity.Customer) error
		Get(id string) (*entity.Customer, error)
		GetByEmail(email string) (*entity.Customer, error)
		GetByUsername(username string) (*entity.Customer, error)
		GetByUsernameOrEmail(username, email string) (*entity.Customer, error)
		List(limit, offset int, filters map[string]interface{}) ([]entity.Customer, int64, error)
		Update(model *entity.Customer) error
		// Log
		CreateLog(model entity.CustomerLog) error
	}
)
