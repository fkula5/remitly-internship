package models

import "gorm.io/gorm"

type BankStock struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Quantity int    `gorm:"not null;check:quantity >= 0" json:"quantity"`
}