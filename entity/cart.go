package entity

import "github.com/google/uuid"

type Cart struct {
	ID         uuid.UUID     `gorm:"primaryKey;type:varchar(36)"`
	Products   []CartProduct `gorm:"type:json;serializer:json;default:null"`
	Total      int
	CustomerID uuid.UUID `gorm:"constraint:OnDelete:CASCADE;"`
}

type CartProduct struct {
	ProductID string `json:"product_id"`
	Qty       int    `json:"qty"`
	Notes     string `json:"notes"`
}
