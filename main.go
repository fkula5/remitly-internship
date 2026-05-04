package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fkula5/remitly-internship/database"
	"github.com/fkula5/remitly-internship/internal/handlers"
	"github.com/fkula5/remitly-internship/internal/repository"
	"github.com/fkula5/remitly-internship/internal/service"
)

func main() {
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	bankRepo := repository.NewBankRepository()
	walletRepo := repository.NewWalletRepository()
	auditRepo := repository.NewAuditLogRepository()
	stockService := service.NewStockService(database.DB, bankRepo, walletRepo, auditRepo)
	stockHandler := handlers.NewStockHandler(stockService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /wallets/{wallet_id}/stocks/{stock_name}", stockHandler.HandleTrade)
	mux.HandleFunc("GET /wallets/{wallet_id}", stockHandler.HandleGetWallet)
	mux.HandleFunc("GET /wallets/{wallet_id}/stocks/{stock_name}", stockHandler.HandleGetWalletStock)
	mux.HandleFunc("GET /stocks", stockHandler.HandleGetBankStocks)
	mux.HandleFunc("POST /stocks", stockHandler.HandleSetBankStocks)
	mux.HandleFunc("GET /log", stockHandler.HandleGetLog)
	mux.HandleFunc("POST /chaos", stockHandler.HandleChaos)

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
