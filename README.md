# Stock Market Simulator

A simplified stock market service with High Availability.

## Requirements
- Docker
- Docker Compose

## Quick Start
To start the application on a specific port (e.g., 8080), run:
```bash
./start.sh 8080
```

This command will:
1. Build and start a PostgreSQL database.
2. Start 2 instances of the Go application for High Availability.
3. Start an Nginx load balancer to distribute traffic across the application instances.

## API Endpoints

### Trade
`POST /wallets/{wallet_id}/stocks/{stock_name}`
Body: `{"type": "buy|sell"}`

### Get Wallet
`GET /wallets/{wallet_id}`

### Get Wallet Stock Quantity
`GET /wallets/{wallet_id}/stocks/{stock_name}`

### Bank State
`GET /stocks`
`POST /stocks` (Set state)

### Audit Log
`GET /log`

### Chaos
`POST /chaos` (Kills the instance handling the request)
