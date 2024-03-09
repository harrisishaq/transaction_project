package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `gorm:"primaryKey;type:binary(16)"`
	Name          string    `gorm:"size:255;" json:"name"`
	Email         string    `gorm:"size:255;unique" json:"email"`
	Password      string
	LastLoginDate *time.Time `gorm:"default:null"`
	Audit         *Audit     `gorm:"type:json;serializer:json;default:null"`
}

type UserLog struct {
	ID            string `gorm:"primaryKey;size:50"`
	Name          string `gorm:"size:255;" json:"name"`
	Email         string `gorm:"size:255;unique" json:"email"`
	Password      string
	LastLoginDate *time.Time `gorm:"default:null"`
	Audit         *Audit     `gorm:"type:json;serializer:json;default:null"`
}
