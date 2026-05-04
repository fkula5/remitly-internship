package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/fkula5/remitly-internship/internal/models"
	"github.com/fkula5/remitly-internship/internal/service"
)

type StockHandler struct {
	stockService *service.StockService
}

func NewStockHandler(ss *service.StockService) *StockHandler {
	return &StockHandler{stockService: ss}
}

func (h *StockHandler) HandleTrade(w http.ResponseWriter, r *http.Request) {
	walletID := r.PathValue("wallet_id")
	stockName := r.PathValue("stock_name")

	var body struct {
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.stockService.Trade(walletID, stockName, body.Type)
	if err != nil {
		if errors.Is(err, service.ErrStockNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, service.ErrInsufficientStock) || errors.Is(err, service.ErrInsufficientWallet) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *StockHandler) HandleGetWallet(w http.ResponseWriter, r *http.Request) {
	walletID := r.PathValue("wallet_id")

	wallet, err := h.stockService.GetWallet(walletID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if wallet == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":     walletID,
			"stocks": []interface{}{},
		})
		return
	}

	response := struct {
		ID     string `json:"id"`
		Stocks []struct {
			Name     string `json:"name"`
			Quantity int    `json:"quantity"`
		} `json:"stocks"`
	}{
		ID: wallet.ExternalID,
	}

	for _, ws := range wallet.Stocks {
		response.Stocks = append(response.Stocks, struct {
			Name     string `json:"name"`
			Quantity int    `json:"quantity"`
		}{
			Name:     ws.Stock.Name,
			Quantity: ws.Quantity,
		})
	}

	json.NewEncoder(w).Encode(response)
}

func (h *StockHandler) HandleGetWalletStock(w http.ResponseWriter, r *http.Request) {
	walletID := r.PathValue("wallet_id")
	stockName := r.PathValue("stock_name")

	quantity, err := h.stockService.GetWalletStockQuantity(walletID, stockName)
	if err != nil {
		if errors.Is(err, service.ErrStockNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(quantity)
}

func (h *StockHandler) HandleGetBankStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := h.stockService.GetBankStocks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Stocks []models.BankStock `json:"stocks"`
	}{
		Stocks: stocks,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *StockHandler) HandleSetBankStocks(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Stocks []models.BankStock `json:"stocks"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.stockService.SetBankStocks(body.Stocks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *StockHandler) HandleGetLog(w http.ResponseWriter, r *http.Request) {
	logs, err := h.stockService.GetAuditLogs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Log []models.AuditLog `json:"log"`
	}{
		Log: logs,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *StockHandler) HandleChaos(w http.ResponseWriter, r *http.Request) {
	os.Exit(1)
}
