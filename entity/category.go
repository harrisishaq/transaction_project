package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        uuid.UUID `gorm:"primaryKey;type:varchar(36)"`
	Name      string
	IsActive  bool
	Products  []Product
	Audit     *Audit         `gorm:"type:json;serializer:json;default:null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type CategoryLog struct {
	ID        string `gorm:"size:50"`
	Name      string
	IsActive  bool
	Audit     *Audit      `gorm:"type:json;serializer:json;default:null"`
	DeletedAt interface{} `gorm:"type:json"`
}
