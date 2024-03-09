package model

import "time"

// request
type (
	CreateUserRequest struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	UpdateUserRequest struct {
		ID    string `json:"id"`
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	ListUserRequest struct {
		Page   int                    `json:"page"`
		Limit  int                    `json:"limit"`
		Filter map[string]interface{} `json:"filters"`
	}
)

// response
type (
	DataUserResponse struct {
		ID            string     `json:"id"`
		Name          string     `json:"name"`
		Email         string     `json:"email"`
		LastLoginDate *time.Time `json:"lastLoginDate"`
	}
)
