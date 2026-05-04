package models

import "gorm.io/gorm"

type AuditLog struct {
	gorm.Model
	Type      string `gorm:"type:varchar(10);not null" json:"type"`
	WalletID  string `gorm:"type:varchar(255);not null" json:"wallet_id"`
	StockName string `gorm:"type:varchar(255);not null" json:"stock_name"`
}

func (AuditLog) TableName() string {
	return "audit_log"
}