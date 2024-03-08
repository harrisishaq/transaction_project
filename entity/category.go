package entity

import "github.com/google/uuid"

type Category struct {
	ID       uuid.UUID `gorm:"primaryKey;type:binary(16)"`
	Name     string
	IsActive bool
	Products []Product
	Audit    *Audit `gorm:"type:json;serializer:json;default:null"`
}

type CategoryLog struct {
	ID       string `gorm:"type:varchar;size:50"`
	Name     string
	IsActive bool
	Products []Product
	Audit    *Audit `gorm:"type:json;serializer:json;default:null"`
}
