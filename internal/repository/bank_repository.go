package repository

import (
	"github.com/fkula5/remitly-internship/internal/models"
	"gorm.io/gorm"
)

type BankRepository interface {
	GetAllStocks(db *gorm.DB) ([]models.BankStock, error)
	GetStockByName(db *gorm.DB, name string) (*models.BankStock, error)
	UpdateStock(db *gorm.DB, stock *models.BankStock) error
	BatchUpsertStocks(db *gorm.DB, stocks []models.BankStock) error
}

type gormBankRepository struct{}

func NewBankRepository() BankRepository {
	return &gormBankRepository{}
}

func (r *gormBankRepository) GetAllStocks(db *gorm.DB) ([]models.BankStock, error) {
	var stocks []models.BankStock
	err := db.Find(&stocks).Error
	return stocks, err
}

func (r *gormBankRepository) GetStockByName(db *gorm.DB, name string) (*models.BankStock, error) {
	var stock models.BankStock
	err := db.Where("name = ?", name).First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *gormBankRepository) UpdateStock(db *gorm.DB, stock *models.BankStock) error {
	return db.Save(stock).Error
}

func (r *gormBankRepository) BatchUpsertStocks(db *gorm.DB, stocks []models.BankStock) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.BankStock{}).Error; err != nil {
			return err
		}
		if len(stocks) > 0 {
			return tx.Create(&stocks).Error
		}
		return nil
	})
}
