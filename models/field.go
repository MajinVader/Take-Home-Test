package models

import "gorm.io/gorm"

type Field struct {
	gorm.Model
	Name         string `json:"name"`
	PricePerHour int    `json:"price_per_hour"`
	Location     string `json:"location"`
}
