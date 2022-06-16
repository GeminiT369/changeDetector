package dto

import "gorm.io/gorm"

type WebItem struct {
	gorm.Model
	Name string `json:"name" gorm:"index"`
	Type string `json:"type" gorm:"index"`
	Url  string `json:"url"`
	Md5  string `json:"md5"`
}
