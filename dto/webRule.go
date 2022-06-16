package dto

import "gorm.io/gorm"

type WebRule struct {
	gorm.Model
	Name  string `json:"name" gorm:"uniqueIndex"`
	Md5   string `json:"md5"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Link  string `json:"link"`
}
