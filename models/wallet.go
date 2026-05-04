package models

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	ExternalID string        `gorm:"type:varchar(255);uniqueIndex;not null" json:"id"`
	Stocks     []WalletStock `json:"stocks,omitempty"`
}