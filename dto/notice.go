package dto

import "gorm.io/gorm"

type Notice struct {
	gorm.Model
	App   string `json:"app" gorm:"index"`
	Mark  uint8  `json:"mark" gorm:"default:0"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Link  string `json:"link"`
}
