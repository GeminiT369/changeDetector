package dto

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `json:"userName" gorm:"index" yaml:"userName"`
	PassWord string `json:"passWord" gorm:"not null" yaml:"passWord"`
}
