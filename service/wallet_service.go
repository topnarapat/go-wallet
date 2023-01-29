package service

import (
	"database/sql"

	"github.com/topnarapat/go-wallet/errs"
	"github.com/topnarapat/go-wallet/logs"
	"github.com/topnarapat/go-wallet/repository"
)

type walletService struct {
	walletRepo repository.WalletRepository
}

func NewWalletService(walletRepo repository.WalletRepository) WalletService {
	return walletService{walletRepo: walletRepo}
}

func (s walletService) ListAllWallets() ([]WalletResponse, error) {
	wallets, err := s.walletRepo.GetAllWallets()
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	walletResponses := []WalletResponse{}
	for _, wallet := range wallets {
		walletResponse := WalletResponse{
			WalletID:  wallet.WalletID,
			Balance:   wallet.Balance,
			Status:    wallet.Status,
			CreatedAt: wallet.CreatedAt,
		}
		walletResponses = append(walletResponses, walletResponse)
	}

	return walletResponses, nil
}

func (s walletService) GetWalletDetail(id int64) (*WalletResponse, error) {
	wallet, err := s.walletRepo.GetWallet(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("wallet not found")
		}

		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	walletResponse := WalletResponse{
		WalletID:  wallet.WalletID,
		Balance:   wallet.Balance,
		Status:    wallet.Status,
		CreatedAt: wallet.CreatedAt,
	}

	return &walletResponse, nil
}

func (s walletService) CreateWallet(w WalletRequest) (*WalletResponse, error) {
	wallet, err := s.walletRepo.CreateNewWallet(w.Balance)
	if err != nil {
		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	walletResponse := WalletResponse{
		WalletID:  wallet.WalletID,
		Balance:   wallet.Balance,
		Status:    wallet.Status,
		CreatedAt: wallet.CreatedAt,
	}

	return &walletResponse, nil
}

func (s walletService) SetWalletBalance(id int64, w AddWalletRequest) (*WalletResponse, error) {
	wallet, err := s.walletRepo.GetWallet(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("wallet not found")
		}

		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	if w.Operation == "Deduct" && wallet.Balance < w.Balance {
		return nil, errs.NewBadRequest("balance not enough")
	}

	if w.Operation == "Deduct" && wallet.Balance >= w.Balance {
		wallet, err = s.walletRepo.SetBalance(id, -w.Balance)
		if err != nil {
			return nil, errs.NewUnexpectedError()
		}
	}

	if w.Operation == "Add" {
		wallet, err = s.walletRepo.SetBalance(id, w.Balance)
		if err != nil {
			return nil, errs.NewUnexpectedError()
		}
	}

	walletResponse := WalletResponse{
		WalletID:  wallet.WalletID,
		Balance:   wallet.Balance,
		Status:    wallet.Status,
		CreatedAt: wallet.CreatedAt,
	}

	return &walletResponse, nil
}

func (s walletService) SetStatusWallet(id int64, st StatusWalletRequest) (*WalletResponse, error) {
	wallet, err := s.walletRepo.SetStatusWallet(id, st.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("wallet not found")
		}

		logs.Error(err)
		return nil, errs.NewUnexpectedError()
	}

	walletResponse := WalletResponse{
		WalletID:  wallet.WalletID,
		Balance:   wallet.Balance,
		Status:    wallet.Status,
		CreatedAt: wallet.CreatedAt,
	}

	return &walletResponse, nil
}
