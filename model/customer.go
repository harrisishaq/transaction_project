package model

import "time"

// request
type (
	CreateCustomerRequest struct {
		Name        string `json:"name" validate:"required"`
		Email       string `json:"email" validate:"required,email"`
		Password    string `json:"password" validate:"required,min=8"`
		Username    string `json:"username" validate:"required,min=8"`
		PhoneNumber string `json:"phone" validate:"required,max=13"`
	}

	ListCustomerRequest struct {
		Page   int                    `json:"page"`
		Limit  int                    `json:"limit"`
		Filter map[string]interface{} `json:"filters"`
	}

	LoginCustomerRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password" validate:"required"`
	}

	UpdateCustomerRequest struct {
		ID    string `json:"id"`
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	UpdateSessionCustomerRequest struct {
		ID      string `json:"id"`
		Session string `json:"session"`
	}
)

// response
type (
	DataCustomerResponse struct {
		ID            string     `json:"id"`
		Name          string     `json:"name"`
		Email         string     `json:"email"`
		Username      string     `json:"username"`
		PhoneNumber   string     `json:"phoneNumber"`
		LastLoginDate *time.Time `json:"lastLoginDate"`
		Session       string     `json:"session"`
	}
)
