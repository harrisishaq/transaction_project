package model

type (
	AddItemCartRequest struct {
		ProductID string `json:"product_id" validate:"required"`
		Qty       int    `json:"qty"`
		Notes     string `json:"notes"`
	}
)

type (
	GetCartResponse struct {
		ID         string      `json:"id"`
		Products   CartProduct `json:"products"`
		Total      int         `json:"total"`
		CustomerID string      `json:"customerId"`
	}

	CartProduct struct {
		AvailableProduct   []Product `json:"availableProduct"`
		UnavailableProduct []Product `json:"unavailableProduct"`
	}

	Product struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Qty         int    `json:"qty"`
		Price       int    `json:"price"`
		PriceTotal  int    `json:"priceTotal"`
		Category    string `json:"category"`
		Notes       string `json:"notes"`
		SystemNotes string `json:"systemNotes"`
	}
)
