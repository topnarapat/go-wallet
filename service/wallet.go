package service

import "time"

type WalletRequest struct {
	Balance float64 `json:"balance"`
}

type AddWalletRequest struct {
	Balance   float64 `json:"balance"`
	Operation string  `json:"operation"`
}

type StatusWalletRequest struct {
	Status string `json:"status"`
}

type WalletResponse struct {
	WalletID  int64     `json:"wallet_id"`
	Balance   float64   `json:"balance"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type WalletService interface {
	ListAllWallets() ([]WalletResponse, error)
	GetWalletDetail(int64) (*WalletResponse, error)
	CreateWallet(WalletRequest) (*WalletResponse, error)
	SetWalletBalance(int64, AddWalletRequest) (*WalletResponse, error)
	SetStatusWallet(int64, StatusWalletRequest) (*WalletResponse, error)
}
