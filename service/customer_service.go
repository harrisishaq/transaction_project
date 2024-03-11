package service

import (
	"test_project/entity"
	"test_project/model"
)

type (
	CustomerService interface {
		CreateCustomer(req *model.CreateCustomerRequest) error
		DeleteCustomer(id string) error
		GenerateTokenAndSession(dataCustomer entity.Customer) (string, error)
		GetCustomer(id string) (*model.DataCustomerResponse, error)
		ListCustomer(req model.ListCustomerRequest) ([]model.DataCustomerResponse, int64, error)
		LoginCustomer(req *model.LoginCustomerRequest) (string, error)
		UpdateCustomer(req *model.UpdateCustomerRequest) error
	}
)
