package entity

import "github.com/google/uuid"

type Category struct {
	ID       uuid.UUID `gorm:"primaryKey;type:binary(16)"`
	Name     string
	IsActive bool
	Products []Product
	Audit    *Audit `gorm:"type:json;serializer:json;default:null"`
}
