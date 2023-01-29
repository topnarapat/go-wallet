package repository

import "database/sql"

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) WalletRepository {
	return walletRepository{db: db}
}

func (r walletRepository) GetAllWallets() ([]Wallet, error) {
	rows, err := r.db.Query("SELECT wallet_id, balance, wallet_status, created_at FROM wallets")
	if err != nil {
		return nil, err
	}

	wallets := []Wallet{}
	for rows.Next() {
		w := Wallet{}
		err = rows.Scan(&w.WalletID, &w.Balance, &w.Status, &w.CreatedAt)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, w)
	}

	return wallets, nil
}

func (r walletRepository) GetWallet(id int64) (*Wallet, error) {
	row := r.db.QueryRow("SELECT wallet_id, balance, wallet_status, created_at FROM wallets WHERE wallet_id=$1", id)
	wallet := Wallet{}
	err := row.Scan(&wallet.WalletID, &wallet.Balance, &wallet.Status, &wallet.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r walletRepository) CreateNewWallet(amount float64) (*Wallet, error) {
	row := r.db.QueryRow("INSERT INTO wallets (balance) values ($1) RETURNING wallet_id, balance, wallet_status, created_at", amount)
	wallet := Wallet{}
	err := row.Scan(&wallet.WalletID, &wallet.Balance, &wallet.Status, &wallet.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r walletRepository) SetBalance(id int64, balance float64) (*Wallet, error) {
	row := r.db.QueryRow("UPDATE wallets SET balance=balance+$2 WHERE wallet_id=$1 RETURNING wallet_id, balance, wallet_status, created_at", id, balance)
	wallet := Wallet{}
	err := row.Scan(&wallet.WalletID, &wallet.Balance, &wallet.Status, &wallet.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r walletRepository) SetStatusWallet(id int64, status string) (*Wallet, error) {
	row := r.db.QueryRow("UPDATE wallets SET wallet_status=$2 WHERE wallet_id=$1 RETURNING wallet_id, balance, wallet_status, created_at", id, status)
	wallet := Wallet{}
	err := row.Scan(&wallet.WalletID, &wallet.Balance, &wallet.Status, &wallet.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}
