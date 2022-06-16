package dto

import "gorm.io/gorm"

type CronTask struct {
	gorm.Model
	TaskName string `json:"taskName" gorm:"not null" yaml:"taskName"`
	CronTime string `json:"cronTime" gorm:"not null" yaml:"cronTime"`
}
