package entity

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `gorm:"primaryKey;type:varchar(36)"`
	Name        string
	Qty         int
	Price       int
	Description string
	IsActive    bool
	CategoryID  uuid.UUID
	Category    Category `gorm:"foreignKey:CategoryID"`
	Audit       *Audit   `gorm:"type:json;serializer:json;default:null"`
}

type ProductLog struct {
	ID          string `gorm:"size:50"`
	Name        string
	Qty         int
	Price       int
	Description string
	IsActive    bool
	CategoryID  uuid.UUID
	Audit       *Audit `gorm:"type:json;serializer:json;default:null"`
}
