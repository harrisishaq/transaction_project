package entity

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID              uuid.UUID `gorm:"primaryKey;type:varchar(36)"`
	Name            string    `gorm:"size:255;" json:"name"`
	Username        string    `gorm:"size:255;unique" json:"username"`
	Email           string    `gorm:"size:255;unique" json:"email"`
	Password        string
	ShippingAddress []ShippingAddress `gorm:"type:json;serializer:json;default:null"`
	PhoneNumber     string
	LastLoginDate   *time.Time `gorm:"default:null"`
	Session         string
	Audit           *Audit `gorm:"type:json;serializer:json;default:null"`
}

type ShippingAddress struct {
	ID              uuid.UUID `gorm:"type:varchar(36)" json:"id"`
	LabelAddress    string    `gorm:"type:varchar(20)" json:"label_address"`
	Address         string    `json:"address"`
	ReceiverName    string    `gorm:"type:varchar(20)" json:"receiver_name"`
	ReceiverContact string    `gorm:"type:varchar(13)" json:"receiver_contact"`
	Notes           string    `json:"notes"`
	IsMain          bool      `json:"is_main"`
}

type CustomerLog struct {
	ID              string `gorm:"primaryKey;size:50"`
	Name            string `gorm:"size:255;" json:"name"`
	Username        string `gorm:"size:255;" json:"username"`
	Email           string `gorm:"size:255;" json:"email"`
	Password        string
	ShippingAddress interface{} `gorm:"type:json"`
	PhoneNumber     string
	LastLoginDate   *time.Time `gorm:"default:null"`
	Session         string
	Audit           *Audit `gorm:"type:json;serializer:json;default:null"`
}
