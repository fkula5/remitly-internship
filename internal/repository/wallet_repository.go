package repository

import (
	"errors"

	"github.com/fkula5/remitly-internship/internal/models"
	"gorm.io/gorm"
)

type WalletRepository interface {
	GetByID(db *gorm.DB, externalID string) (*models.Wallet, error)
	Create(db *gorm.DB, wallet *models.Wallet) error
	GetStock(db *gorm.DB, walletID uint, stockID uint) (*models.WalletStock, error)
	UpdateStock(db *gorm.DB, stock *models.WalletStock) error
	CreateStock(db *gorm.DB, stock *models.WalletStock) error
}

type gormWalletRepository struct{}

func NewWalletRepository() WalletRepository {
	return &gormWalletRepository{}
}

func (r *gormWalletRepository) GetByID(db *gorm.DB, externalID string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := db.Preload("Stocks.Stock").Where("external_id = ?", externalID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *gormWalletRepository) Create(db *gorm.DB, wallet *models.Wallet) error {
	return db.Create(wallet).Error
}

func (r *gormWalletRepository) GetStock(db *gorm.DB, walletID uint, stockID uint) (*models.WalletStock, error) {
	var ws models.WalletStock
	err := db.Where("wallet_id = ? AND stock_id = ?", walletID, stockID).First(&ws).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ws, nil
}

func (r *gormWalletRepository) UpdateStock(db *gorm.DB, stock *models.WalletStock) error {
	return db.Save(stock).Error
}

func (r *gormWalletRepository) CreateStock(db *gorm.DB, stock *models.WalletStock) error {
	return db.Create(stock).Error
}
