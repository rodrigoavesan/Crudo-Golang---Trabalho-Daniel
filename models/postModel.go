package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	ID     string `gorm:"primaryKey"`
	Title  string
	Artist string
	Price  float64
	Gender string
}
