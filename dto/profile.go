package dto

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	SmtpUser   string `json:"smtpUser" gorm:"index" yaml:"smtpUser"`
	SmtpPswd   string `json:"smtpPswd" gorm:"not null" yaml:"smtpPswd"`
	SmtpServer string `json:"smtpServer" gorm:"not null" yaml:"smtpServer"`
	SmtpPort   string `json:"smtpPort" gorm:"not null" yaml:"smtpPort"`
}
