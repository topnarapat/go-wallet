package repository

import "time"

type WalletRepository interface {
	GetAllWallets() ([]Wallet, error)
	GetWallet(int64) (*Wallet, error)
	CreateNewWallet(float64) (*Wallet, error)
	SetBalance(int64, float64) (*Wallet, error)
	SetStatusWallet(int64, string) (*Wallet, error)
}

type Wallet struct {
	WalletID  int64     `db:"wallet_id"`
	Balance   float64   `db:"balance"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
}
