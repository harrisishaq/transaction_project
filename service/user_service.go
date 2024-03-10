package service

import (
	"test_project/entity"
	"test_project/model"
)

type (
	UserService interface {
		CreateUser(req *model.CreateUserRequest) error
		DeleteUser(id string) error
		GenerateTokenAndSession(dataUser entity.User) (string, error)
		GetUser(id string) (*model.DataUserResponse, error)
		ListUser(req model.ListUserRequest) ([]model.DataUserResponse, int64, error)
		LoginUser(req *model.LoginUserRequest) (string, error)
		UpdateUser(req *model.UpdateUserRequest) error
	}
)
