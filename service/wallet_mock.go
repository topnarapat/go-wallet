package service

import "github.com/stretchr/testify/mock"

type walletServiceMock struct {
	mock.Mock
}

func NewWalletServiceMock() *walletServiceMock {
	return &walletServiceMock{}
}

func (s *walletServiceMock) ListAllWallets() ([]WalletResponse, error) {
	args := s.Called()
	return args.Get(0).([]WalletResponse), args.Error(1)
}

func (s *walletServiceMock) GetWalletDetail(id int64) (*WalletResponse, error) {
	args := s.Called(id)
	return args.Get(0).(*WalletResponse), args.Error(1)
}

func (s *walletServiceMock) CreateWallet(r WalletRequest) (*WalletResponse, error) {
	args := s.Called(r)
	return args.Get(0).(*WalletResponse), args.Error(1)
}

func (s *walletServiceMock) SetWalletBalance(id int64, r AddWalletRequest) (*WalletResponse, error) {
	args := s.Called(id, r)
	return args.Get(0).(*WalletResponse), args.Error(1)
}

func (s *walletServiceMock) SetStatusWallet(id int64, r StatusWalletRequest) (*WalletResponse, error) {
	args := s.Called(id, r)
	return args.Get(0).(*WalletResponse), args.Error(1)
}
