package service

import (
	"context"
	"test_project/entity"
	"test_project/model"
)

type (
	UserService interface {
		CreateUser(ctx context.Context, req *model.CreateUserRequest) error
		DeleteUser(ctx context.Context, id string) error
		GenerateTokenAndSession(dataUser entity.User) (string, error)
		GetUser(id string) (*model.DataUserResponse, error)
		ListUser(req model.ListUserRequest) ([]model.DataUserResponse, int64, error)
		LoginUser(req *model.LoginUserRequest) (string, error)
		LogoutUser(ctx context.Context, id string) error
		UpdateSesionUser(req *model.UpdateSessionUserRequest) error
		UpdateUser(ctx context.Context, req *model.UpdateUserRequest) error
	}
)
