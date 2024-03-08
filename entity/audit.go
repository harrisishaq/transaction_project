package entity

import "time"

type Audit struct {
	LogReason *string    `json:"log_reason,omitempty"`
	CurrNo    int        `json:"curr_no"`
	CreatedAt *time.Time `json:"created_at"`
	CreatedBy string     `json:"created_by"`
	UpdatedAt *time.Time `gorm:"default:null" json:"updated_at"`
	UpdatedBy string     `gorm:"default:null" json:"updated_by"`
}
