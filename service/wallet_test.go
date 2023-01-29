//go:build unit
// +build unit

package service_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/topnarapat/go-wallet/errs"
	"github.com/topnarapat/go-wallet/repository"
	"github.com/topnarapat/go-wallet/service"
)

func TestListAllWallets(t *testing.T) {
	t.Run("get all wallets", func(t *testing.T) {
		// Arrange
		walletRepo := repository.NewWalletRepositoryMock()
		walletRepo.On("GetAllWallets").Return([]repository.Wallet{
			{WalletID: 1, Balance: 500, Status: "Active", CreatedAt: time.Date(2022, time.January, 27, 12, 30, 0, 0, time.UTC)},
			{WalletID: 2, Balance: 0, Status: "Deactive", CreatedAt: time.Date(2022, time.January, 27, 12, 30, 0, 0, time.UTC)},
		}, nil)

		walletService := service.NewWalletService(walletRepo)

		// Act
		wallets, _ := walletService.ListAllWallets()
		expected := []service.WalletResponse{
			{WalletID: 1, Balance: 500, Status: "Active", CreatedAt: time.Date(2022, time.January, 27, 12, 30, 0, 0, time.UTC)},
			{WalletID: 2, Balance: 0, Status: "Deactive", CreatedAt: time.Date(2022, time.January, 27, 12, 30, 0, 0, time.UTC)},
		}

		// Assert
		assert.Equal(t, expected, wallets)
	})

	t.Run("unexpected error", func(t *testing.T) {
		// Arrange
		walletRepo := repository.NewWalletRepositoryMock()
		walletRepo.On("GetAllWallets").Return([]repository.Wallet{}, errors.New(""))

		walletService := service.NewWalletService(walletRepo)

		// Act
		_, err := walletService.ListAllWallets()

		// Assert
		assert.ErrorIs(t, err, errs.NewUnexpectedError())
	})
}

func TestGetWalletDetail(t *testing.T) {
	type testCase struct {
		name      string
		walletID  int64
		balance   float64
		status    string
		createdAt time.Time
	}

	cases := []testCase{
		{name: "get wallet id 1", walletID: 1, balance: 100, status: "Active", createdAt: time.Date(2022, time.January, 27, 12, 30, 0, 0, time.UTC)},
		{name: "get wallet id 2", walletID: 2, balance: 200, status: "Active", createdAt: time.Date(2022, time.January, 28, 12, 30, 0, 0, time.UTC)},
		{name: "get wallet id 3", walletID: 3, balance: 300, status: "Deactive", createdAt: time.Date(2022, time.January, 29, 12, 30, 0, 0, time.UTC)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Arrange
			walletRepo := repository.NewWalletRepositoryMock()
			walletRepo.On("GetWallet", c.walletID).Return(&repository.Wallet{
				WalletID:  c.walletID,
				Balance:   c.balance,
				Status:    c.status,
				CreatedAt: c.createdAt,
			}, nil)

			walletService := service.NewWalletService(walletRepo)

			// Act
			wallet, _ := walletService.GetWalletDetail(c.walletID)
			expected := &service.WalletResponse{
				WalletID:  c.walletID,
				Balance:   c.balance,
				Status:    c.status,
				CreatedAt: c.createdAt,
			}

			// Assert
			assert.Equal(t, expected, wallet)
		})
	}

	t.Run("wallet not found", func(t *testing.T) {
		// Arrange
		var id int64 = 99
		walletRepo := repository.NewWalletRepositoryMock()
		walletRepo.On("GetWallet", id).Return(&repository.Wallet{}, sql.ErrNoRows)

		walletService := service.NewWalletService(walletRepo)

		// Act
		_, err := walletService.GetWalletDetail(id)

		// Assert
		assert.ErrorIs(t, err, errs.NewNotFoundError("wallet not found"))
	})

	t.Run("unexpected error", func(t *testing.T) {
		// Arrange
		var id int64 = 99
		walletRepo := repository.NewWalletRepositoryMock()
		walletRepo.On("GetWallet", id).Return(&repository.Wallet{}, errors.New(""))

		walletService := service.NewWalletService(walletRepo)

		// Act
		_, err := walletService.GetWalletDetail(id)

		// Assert
		assert.ErrorIs(t, err, errs.NewUnexpectedError())
	})
}

func TestCreateWallet(t *testing.T) {
	type testCase struct {
		name      string
		walletID  int64
		balance   float64
		status    string
		createdAt time.Time
	}

	cases := []testCase{
		{name: "create wallet id 1", walletID: 1, balance: 100, status: "Active", createdAt: time.Date(2022, time.January, 27, 12, 30, 0, 0, time.UTC)},
		{name: "create wallet id 2", walletID: 2, balance: 200, status: "Active", createdAt: time.Date(2022, time.January, 28, 12, 30, 0, 0, time.UTC)},
		{name: "create wallet id 3", walletID: 3, balance: 300, status: "Active", createdAt: time.Date(2022, time.January, 29, 12, 30, 0, 0, time.UTC)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Arrange
			walletRepo := repository.NewWalletRepositoryMock()
			walletRepo.On("CreateNewWallet", c.balance).Return(&repository.Wallet{
				WalletID:  c.walletID,
				Balance:   c.balance,
				Status:    c.status,
				CreatedAt: c.createdAt,
			}, nil)

			walletService := service.NewWalletService(walletRepo)

			balance := service.WalletRequest{
				Balance: c.balance,
			}

			// Act
			wallet, _ := walletService.CreateWallet(balance)
			expected := &service.WalletResponse{
				WalletID:  c.walletID,
				Balance:   c.balance,
				Status:    c.status,
				CreatedAt: c.createdAt,
			}

			// Assert
			assert.Equal(t, expected, wallet)
		})
	}

	t.Run("unexpected error", func(t *testing.T) {
		// Arrange
		var balance float64 = 99
		walletRepo := repository.NewWalletRepositoryMock()
		walletRepo.On("CreateNewWallet", balance).Return(&repository.Wallet{}, errors.New(""))

		walletService := service.NewWalletService(walletRepo)

		walletRequest := service.WalletRequest{
			Balance: balance,
		}

		// Act
		_, err := walletService.CreateWallet(walletRequest)

		// Assert
		assert.ErrorIs(t, err, errs.NewUnexpectedError())
	})
}

func TestSetWalletBalance(t *testing.T) {
	t.Run("add balance", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		amount := service.AddWalletRequest{
			Balance:   1000,
			Operation: "Add",
		}
		walletRepo := repository.NewWalletRepositoryMock()
		walletRepo.On("GetWallet", id).Return(&repository.Wallet{
			WalletID:  id,
			Balance:   2000,
			Status:    "Active",
			CreatedAt: time.Date(2022, time.January, 29, 12, 30, 0, 0, time.UTC),
		}, nil)
		walletRepo.On("SetBalance", id, amount.Balance).Return(&repository.Wallet{
			WalletID:  id,
			Balance:   3000,
			Status:    "Active",
			CreatedAt: time.Date(2022, time.January, 29, 12, 30, 0, 0, time.UTC),
		}, nil)

		walletService := service.NewWalletService(walletRepo)

		// Act
		wallet, _ := walletService.SetWalletBalance(id, amount)
		expected := &service.WalletResponse{
			WalletID:  id,
			Balance:   3000,
			Status:    "Active",
			CreatedAt: time.Date(2022, time.January, 29, 12, 30, 0, 0, time.UTC),
		}

		// Assert
		assert.Equal(t, wallet, expected)
	})

	t.Run("deduct balance", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		amount := service.AddWalletRequest{
			Balance:   1000,
			Operation: "Deduct",
		}
		walletRepo := repository.NewWalletRepositoryMock()
		walletRepo.On("GetWallet", id).Return(&repository.Wallet{
			WalletID:  id,
			Balance:   3000,
			Status:    "Active",
			CreatedAt: time.Date(2022, time.January, 29, 12, 30, 0, 0, time.UTC),
		}, nil)
		walletRepo.On("SetBalance", id, -amount.Balance).Return(&repository.Wallet{
			WalletID:  id,
			Balance:   2000,
			Status:    "Active",
			CreatedAt: time.Date(2022, time.January, 29, 12, 30, 0, 0, time.UTC),
		}, nil)

		walletService := service.NewWalletService(walletRepo)

		// Act
		wallet, _ := walletService.SetWalletBalance(id, amount)
		expected := &service.WalletResponse{
			WalletID:  id,
			Balance:   2000,
			Status:    "Active",
			CreatedAt: time.Date(2022, time.January, 29, 12, 30, 0, 0, time.UTC),
		}

		// Assert
		assert.Equal(t, wallet, expected)
	})

	t.Run("wallet not found", func(t *testing.T) {
		// Arrange
		var id int64 = 99
		amount := service.AddWalletRequest{
			Balance:   1000,
			Operation: "Add",
		}
		walletRepo := repository.NewWalletRepositoryMock()
		walletRepo.On("GetWallet", id).Return(&repository.Wallet{}, sql.ErrNoRows)

		walletService := service.NewWalletService(walletRepo)

		// Act
		_, err := walletService.SetWalletBalance(id, amount)

		// Assert
		assert.ErrorIs(t, err, errs.NewNotFoundError("wallet not found"))
	})

	t.Run("unexpected error", func(t *testing.T) {
		// Arrange
		var id int64 = 99
		amount := service.AddWalletRequest{
			Balance:   1000,
			Operation: "Add",
		}
		walletRepo := repository.NewWalletRepositoryMock()
		walletRepo.On("GetWallet", id).Return(&repository.Wallet{}, errors.New(""))

		walletService := service.NewWalletService(walletRepo)

		// Act
		_, err := walletService.SetWalletBalance(id, amount)

		// Assert
		assert.ErrorIs(t, err, errs.NewUnexpectedError())
	})

	t.Run("balance not enough", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		amount := service.AddWalletRequest{
			Balance:   1000,
			Operation: "Deduct",
		}
		walletRepo := repository.NewWalletRepositoryMock()
		walletRepo.On("GetWallet", id).Return(&repository.Wallet{
			WalletID:  id,
			Balance:   500,
			Status:    "Active",
			CreatedAt: time.Date(2022, time.January, 29, 12, 30, 0, 0, time.UTC),
		}, nil)
		walletRepo.On("SetBalance", id, amount.Balance).Return(&repository.Wallet{}, errors.New("balance not enough"))

		walletService := service.NewWalletService(walletRepo)

		// Act
		_, err := walletService.SetWalletBalance(id, amount)

		// Assert
		assert.ErrorIs(t, err, errs.NewBadRequest("balance not enough"))
	})

	t.Run("deduct unexpected error", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		amount := service.AddWalletRequest{
			Balance:   1000,
			Operation: "Deduct",
		}
		walletRepo := repository.NewWalletRepositoryMock()
		walletRepo.On("GetWallet", id).Return(&repository.Wallet{
			WalletID:  id,
			Balance:   2000,
			Status:    "Active",
			CreatedAt: time.Date(2022, time.January, 29, 12, 30, 0, 0, time.UTC),
		}, nil)
		walletRepo.On("SetBalance", id, -amount.Balance).Return(&repository.Wallet{}, errors.New(""))

		walletService := service.NewWalletService(walletRepo)

		// Act
		_, err := walletService.SetWalletBalance(id, amount)

		// Assert
		assert.ErrorIs(t, err, errs.NewUnexpectedError())
	})

	t.Run("add unexpected error", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		amount := service.AddWalletRequest{
			Balance:   1000,
			Operation: "Add",
		}
		walletRepo := repository.NewWalletRepositoryMock()
		walletRepo.On("GetWallet", id).Return(&repository.Wallet{
			WalletID:  id,
			Balance:   2000,
			Status:    "Active",
			CreatedAt: time.Date(2022, time.January, 29, 12, 30, 0, 0, time.UTC),
		}, nil)
		walletRepo.On("SetBalance", id, amount.Balance).Return(&repository.Wallet{}, errors.New(""))

		walletService := service.NewWalletService(walletRepo)

		// Act
		_, err := walletService.SetWalletBalance(id, amount)

		// Assert
		assert.ErrorIs(t, err, errs.NewUnexpectedError())
	})
}

func TestSetStatusWallet(t *testing.T) {
	t.Run("set status", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		st := service.StatusWalletRequest{
			Status: "Deactive",
		}
		walletRepo := repository.NewWalletRepositoryMock()
		walletRepo.On("SetStatusWallet", id, st.Status).Return(&repository.Wallet{
			WalletID:  id,
			Balance:   2000,
			Status:    "Deactive",
			CreatedAt: time.Date(2022, time.January, 29, 12, 30, 0, 0, time.UTC),
		}, nil)

		walletService := service.NewWalletService(walletRepo)

		// Act
		wallet, _ := walletService.SetStatusWallet(id, st)
		expected := &service.WalletResponse{
			WalletID:  id,
			Balance:   2000,
			Status:    "Deactive",
			CreatedAt: time.Date(2022, time.January, 29, 12, 30, 0, 0, time.UTC),
		}

		// Assert
		assert.Equal(t, expected, wallet)
	})

	t.Run("wallet not found", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		st := service.StatusWalletRequest{
			Status: "Deactive",
		}
		walletRepo := repository.NewWalletRepositoryMock()
		walletRepo.On("SetStatusWallet", id, st.Status).Return(&repository.Wallet{}, sql.ErrNoRows)

		walletService := service.NewWalletService(walletRepo)

		// Act
		_, err := walletService.SetStatusWallet(id, st)

		// Assert
		assert.ErrorIs(t, err, errs.NewNotFoundError("wallet not found"))
	})

	t.Run("unexpected error", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		st := service.StatusWalletRequest{
			Status: "Deactive",
		}
		walletRepo := repository.NewWalletRepositoryMock()
		walletRepo.On("SetStatusWallet", id, st.Status).Return(&repository.Wallet{}, errors.New(""))

		walletService := service.NewWalletService(walletRepo)

		// Act
		_, err := walletService.SetStatusWallet(id, st)

		// Assert
		assert.ErrorIs(t, err, errs.NewUnexpectedError())
	})
}
