package service

import (
	"errors"
	"fmt"

	"github.com/fkula5/remitly-internship/internal/models"
	"github.com/fkula5/remitly-internship/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrStockNotFound      = errors.New("stock not found")
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrInsufficientWallet = errors.New("insufficient stock in wallet")
)

type StockService struct {
	db        *gorm.DB
	bankRepo  repository.BankRepository
	walletRepo repository.WalletRepository
	auditRepo  repository.AuditLogRepository
}

func NewStockService(db *gorm.DB, br repository.BankRepository, wr repository.WalletRepository, ar repository.AuditLogRepository) *StockService {
	return &StockService{
		db:         db,
		bankRepo:   br,
		walletRepo: wr,
		auditRepo:  ar,
	}
}

func (s *StockService) Trade(walletExternalID, stockName, tradeType string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get or create wallet
		wallet, err := s.walletRepo.GetByID(tx, walletExternalID)
		if err != nil {
			return err
		}
		if wallet == nil {
			wallet = &models.Wallet{ExternalID: walletExternalID}
			if err := s.walletRepo.Create(tx, wallet); err != nil {
				return err
			}
		}

		// 2. Get stock from bank
		bankStock, err := s.bankRepo.GetStockByName(tx, stockName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrStockNotFound
			}
			return err
		}

		// 3. Process trade
		if tradeType == "buy" {
			if bankStock.Quantity < 1 {
				return ErrInsufficientStock
			}
			bankStock.Quantity--
			if err := s.bankRepo.UpdateStock(tx, bankStock); err != nil {
				return err
			}

			ws, err := s.walletRepo.GetStock(tx, wallet.ID, bankStock.ID)
			if err != nil {
				return err
			}
			if ws == nil {
				ws = &models.WalletStock{WalletID: wallet.ID, StockID: bankStock.ID, Quantity: 1}
				if err := s.walletRepo.CreateStock(tx, ws); err != nil {
					return err
				}
			} else {
				ws.Quantity++
				if err := s.walletRepo.UpdateStock(tx, ws); err != nil {
					return err
				}
			}
		} else if tradeType == "sell" {
			ws, err := s.walletRepo.GetStock(tx, wallet.ID, bankStock.ID)
			if err != nil {
				return err
			}
			if ws == nil || ws.Quantity < 1 {
				return ErrInsufficientWallet
			}
			ws.Quantity--
			if err := s.walletRepo.UpdateStock(tx, ws); err != nil {
				return err
			}

			bankStock.Quantity++
			if err := s.bankRepo.UpdateStock(tx, bankStock); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("invalid trade type: %s", tradeType)
		}

		// 4. Audit log (only successful operations)
		audit := &models.AuditLog{
			Type:      tradeType,
			WalletID:  walletExternalID,
			StockName: stockName,
		}
		return s.auditRepo.Create(tx, audit)
	})
}

func (s *StockService) GetWallet(walletExternalID string) (*models.Wallet, error) {
	return s.walletRepo.GetByID(s.db, walletExternalID)
}

func (s *StockService) GetWalletStockQuantity(walletExternalID, stockName string) (int, error) {
	wallet, err := s.walletRepo.GetByID(s.db, walletExternalID)
	if err != nil || wallet == nil {
		return 0, err
	}

	bankStock, err := s.bankRepo.GetStockByName(s.db, stockName)
	if err != nil {
		return 0, err
	}

	ws, err := s.walletRepo.GetStock(s.db, wallet.ID, bankStock.ID)
	if err != nil || ws == nil {
		return 0, err
	}

	return ws.Quantity, nil
}

func (s *StockService) GetBankStocks() ([]models.BankStock, error) {
	return s.bankRepo.GetAllStocks(s.db)
}

func (s *StockService) SetBankStocks(stocks []models.BankStock) error {
	return s.bankRepo.BatchUpsertStocks(s.db, stocks)
}

func (s *StockService) GetAuditLogs() ([]models.AuditLog, error) {
	return s.auditRepo.GetAll(s.db)
}
