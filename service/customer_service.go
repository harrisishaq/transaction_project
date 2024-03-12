package service

import (
	"context"
	"test_project/entity"
	"test_project/model"
)

type (
	CustomerService interface {
		CreateCustomer(ctx context.Context, req *model.CreateCustomerRequest) error
		DeleteCustomer(ctx context.Context, id string) error
		GenerateTokenAndSession(dataCustomer entity.Customer) (string, error)
		GetCustomer(ctx context.Context, id string) (*model.DataCustomerResponse, error)
		GetCustomerByID(id string) (*model.DataCustomerResponse, error)
		ListCustomer(ctx context.Context, req model.ListCustomerRequest) ([]model.DataCustomerResponse, int64, error)
		LoginCustomer(req *model.LoginCustomerRequest) (string, error)
		UpdateCustomer(ctx context.Context, req *model.UpdateCustomerRequest) error
		UpdateSesion(req *model.UpdateSessionCustomerRequest) error
	}
)
