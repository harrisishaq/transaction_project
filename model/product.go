package model

type (
	CreateProductRequest struct {
		Name        string `json:"name" validate:"required"`
		Qty         int    `json:"qty"`
		Price       int    `json:"price"`
		Description string `json:"description"`
		CategoryID  string `json:"categoryId" validate:"required"`
	}

	UpdateProductRequest struct {
		ID          string `json:"id"`
		Name        string `json:"name" validate:"required"`
		IsActive    bool   `json:"isActive"`
		Qty         int    `json:"qty"`
		Price       int    `json:"price"`
		Description string `json:"description"`
		CategoryID  string `json:"categoryId" validate:"required"`
	}

	ListProductRequest struct {
		Page   int                    `json:"page"`
		Limit  int                    `json:"limit"`
		Filter map[string]interface{} `json:"filters"`
	}
)

// response
type (
	DataProductResponse struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		IsActive     bool   `json:"isActive"`
		CategoryID   string `json:"categoryId"`
		CategoryName string `json:"categoryName"`
		Qty          int    `json:"qty"`
		Price        int    `json:"price"`
		Description  string `json:"description"`
	}
)
