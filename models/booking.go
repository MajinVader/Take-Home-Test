package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserID  uint      `json:"user_id"`
	FieldID uint      `json:"field_id"`
	Start   time.Time `json:"start_time"`
	End     time.Time `json:"end_time"`
	Status  string    `json:"status"` // "pending", "paid", "cancelled", dll
}
