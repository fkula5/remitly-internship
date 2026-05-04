package models

import "gorm.io/gorm"

type WalletStock struct {
	gorm.Model
	WalletID uint `gorm:"uniqueIndex:idx_wallet_stock;not null" json:"-"`
	StockID  uint `gorm:"uniqueIndex:idx_wallet_stock;not null" json:"-"`

	Wallet   Wallet    `json:"-"`
	Stock    BankStock `json:"-"`
	
	Quantity int       `gorm:"not null;check:quantity >= 0" json:"quantity"`
}