package model

// request
type (
	CreateCategoryRequest struct {
		Name string `json:"name" validate:"required"`
	}

	UpdateCategoryRequest struct {
		ID       string `json:"id"`
		Name     string `json:"name" validate:"required"`
		IsActive bool   `json:"isActive"`
	}

	ListCategoryRequest struct {
		Page   int                    `json:"page"`
		Limit  int                    `json:"limit"`
		Filter map[string]interface{} `json:"filters"`
	}
)

// response
type (
	DataCategoryResponse struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		IsActive bool   `json:"isActive"`
	}
)
