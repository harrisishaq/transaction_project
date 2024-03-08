package entity

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `gorm:"primaryKey;type:binary(16)"`
	Name        string
	Qty         int
	Price       string
	Description string
	IsActive    bool
	CategoryID  uuid.UUID
	Category    Category `gorm:"foreignKey:CategoryID"`
	Audit       *Audit   `gorm:"type:json;serializer:json;default:null"`
}
