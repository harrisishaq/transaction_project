package service

import "test_project/model"

type (
	UserService interface {
		CreateUser(req *model.CreateUserRequest) error
		DeleteUser(id string) error
		GetUser(id string) (*model.DataUserResponse, error)
		ListUser(req model.ListUserRequest) ([]model.DataUserResponse, int64, error)
		UpdateUser(req *model.UpdateUserRequest) error
	}
)
