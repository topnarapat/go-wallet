package repository

import "github.com/stretchr/testify/mock"

type walletRepositoryMock struct {
	mock.Mock
}

func NewWalletRepositoryMock() *walletRepositoryMock {
	return &walletRepositoryMock{}
}

func (r *walletRepositoryMock) GetAllWallets() ([]Wallet, error) {
	args := r.Called()
	return args.Get(0).([]Wallet), args.Error(1)
}

func (r *walletRepositoryMock) GetWallet(id int64) (*Wallet, error) {
	args := r.Called(id)
	return args.Get(0).(*Wallet), args.Error(1)
}

func (r *walletRepositoryMock) CreateNewWallet(amount float64) (*Wallet, error) {
	args := r.Called(amount)
	return args.Get(0).(*Wallet), args.Error(1)
}

func (r *walletRepositoryMock) SetBalance(id int64, amount float64) (*Wallet, error) {
	args := r.Called(id, amount)
	return args.Get(0).(*Wallet), args.Error(1)
}

func (r *walletRepositoryMock) SetStatusWallet(id int64, status string) (*Wallet, error) {
	args := r.Called(id, status)
	return args.Get(0).(*Wallet), args.Error(1)
}
