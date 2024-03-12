package repository

import "test_project/entity"

type (
	CartRepository interface {
		Create(model entity.Cart) error
		Get(custID string) (*entity.Cart, error)
		Update(model *entity.Cart) error
	}
)
