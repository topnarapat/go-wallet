//go:build unit
// +build unit

package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/topnarapat/go-wallet/errs"
	"github.com/topnarapat/go-wallet/handler"
	"github.com/topnarapat/go-wallet/service"
)

func TestListWallets(t *testing.T) {
	t.Run("get all wallet success", func(t *testing.T) {
		// Arrange
		walletService := service.NewWalletServiceMock()
		walletService.On("ListAllWallets").Return([]service.WalletResponse{
			{WalletID: 1, Balance: 500, Status: "Active", CreatedAt: time.Date(2023, time.January, 27, 12, 30, 0, 0, time.UTC)},
		}, nil)

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expected := `[{"wallet_id":1,"balance":500,"status":"Active","created_at":"2023-01-27T12:30:00Z"}]`

		// Assert
		if assert.NoError(t, walletHandler.ListWallets(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("get all wallet error", func(t *testing.T) {
		// Arrange
		walletService := service.NewWalletServiceMock()
		walletService.On("ListAllWallets").Return([]service.WalletResponse{}, errors.New(""))

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assert
		if assert.NoError(t, walletHandler.ListWallets(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestGetWallet(t *testing.T) {
	t.Run("get wallet by id success", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		walletService := service.NewWalletServiceMock()
		walletService.On("GetWalletDetail", id).Return(&service.WalletResponse{
			WalletID:  id,
			Balance:   500,
			Status:    "Active",
			CreatedAt: time.Date(2023, time.January, 27, 12, 30, 0, 0, time.UTC),
		}, nil)

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/wallet/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		expected := `{"wallet_id":1,"balance":500,"status":"Active","created_at":"2023-01-27T12:30:00Z"}`

		// Assert
		if assert.NoError(t, walletHandler.GetWallet(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("id must be number", func(t *testing.T) {
		// Arrange
		walletService := service.NewWalletServiceMock()
		walletService.On("GetWalletDetail", "abc").Return(&service.WalletResponse{}, errors.New("id must be number"))

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/wallet/:id")
		c.SetParamNames("id")
		c.SetParamValues("abc")

		// Assert
		if assert.NoError(t, walletHandler.GetWallet(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("not found wallet", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		walletService := service.NewWalletServiceMock()
		walletService.On("GetWalletDetail", id).Return(&service.WalletResponse{}, errs.NewNotFoundError("not found wallet"))

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/wallet/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// Assert
		if assert.NoError(t, walletHandler.GetWallet(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
}

func TestCreateWallet(t *testing.T) {
	t.Run("create wallet success", func(t *testing.T) {
		// Arrange
		balance := service.WalletRequest{
			Balance: 1000,
		}
		walletService := service.NewWalletServiceMock()
		walletService.On("CreateWallet", balance).Return(&service.WalletResponse{
			WalletID:  1,
			Balance:   1000,
			Status:    "Active",
			CreatedAt: time.Date(2023, time.January, 27, 12, 30, 0, 0, time.UTC),
		}, nil)

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		r := `{"balance":1000}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(r))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expected := `{"wallet_id":1,"balance":1000,"status":"Active","created_at":"2023-01-27T12:30:00Z"}`

		// Assert
		if assert.NoError(t, walletHandler.CreateWallet(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("request body incorrect format", func(t *testing.T) {
		// Arrange
		balance := service.WalletRequest{
			Balance: 1000,
		}
		walletService := service.NewWalletServiceMock()
		walletService.On("CreateWallet", balance).Return(&service.WalletResponse{}, errors.New(""))

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		r := `{"balance":"abc"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(r))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assert
		if assert.NoError(t, walletHandler.CreateWallet(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("create wallet unexpected error", func(t *testing.T) {
		// Arrange
		balance := service.WalletRequest{
			Balance: 1000,
		}
		walletService := service.NewWalletServiceMock()
		walletService.On("CreateWallet", balance).Return(&service.WalletResponse{}, errors.New(""))

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		r := `{"balance":1000}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(r))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Assert
		if assert.NoError(t, walletHandler.CreateWallet(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestAddBalance(t *testing.T) {
	t.Run("add balance success", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		balance := service.AddWalletRequest{
			Balance:   1000,
			Operation: "Add",
		}
		walletService := service.NewWalletServiceMock()
		walletService.On("SetWalletBalance", id, balance).Return(&service.WalletResponse{
			WalletID:  1,
			Balance:   2000,
			Status:    "Active",
			CreatedAt: time.Date(2023, time.January, 27, 12, 30, 0, 0, time.UTC),
		}, nil)

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		r := `{"balance":1000,"operation":"Add"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(r))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/wallet/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		expected := `{"wallet_id":1,"balance":2000,"status":"Active","created_at":"2023-01-27T12:30:00Z"}`

		// Assert
		if assert.NoError(t, walletHandler.AddBalance(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("deduct balance success", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		balance := service.AddWalletRequest{
			Balance:   1000,
			Operation: "Deduct",
		}
		walletService := service.NewWalletServiceMock()
		walletService.On("SetWalletBalance", id, balance).Return(&service.WalletResponse{
			WalletID:  1,
			Balance:   2000,
			Status:    "Active",
			CreatedAt: time.Date(2023, time.January, 27, 12, 30, 0, 0, time.UTC),
		}, nil)

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		r := `{"balance":1000,"operation":"Deduct"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(r))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/wallet/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		expected := `{"wallet_id":1,"balance":2000,"status":"Active","created_at":"2023-01-27T12:30:00Z"}`

		// Assert
		if assert.NoError(t, walletHandler.AddBalance(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("id must be number", func(t *testing.T) {
		// Arrange
		balance := service.AddWalletRequest{
			Balance:   1000,
			Operation: "Deduct",
		}
		walletService := service.NewWalletServiceMock()
		walletService.On("SetWalletBalance", "abc", balance).Return(&service.WalletResponse{}, errors.New("id must be number"))

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		r := `{"balance":1000,"operation":"Deduct"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(r))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/wallet/:id")
		c.SetParamNames("id")
		c.SetParamValues("abc")

		// Assert
		if assert.NoError(t, walletHandler.AddBalance(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("request body incorrect format", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		balance := service.AddWalletRequest{}
		walletService := service.NewWalletServiceMock()
		walletService.On("SetWalletBalance", id, balance).Return(&service.WalletResponse{}, errors.New("request body incorrect format"))

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		r := `{"balance":"1000","operation":"Deduct"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(r))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/wallet/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// Assert
		if assert.NoError(t, walletHandler.AddBalance(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("add balance unexpected error", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		balance := service.AddWalletRequest{
			Balance:   1000,
			Operation: "Add",
		}
		walletService := service.NewWalletServiceMock()
		walletService.On("SetWalletBalance", id, balance).Return(&service.WalletResponse{}, errors.New(""))

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		r := `{"balance":1000,"operation":"Add"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(r))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/wallet/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// Assert
		if assert.NoError(t, walletHandler.AddBalance(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}

func TestChangeStatus(t *testing.T) {
	t.Run("change status success", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		request := service.StatusWalletRequest{
			Status: "Deactive",
		}
		walletService := service.NewWalletServiceMock()
		walletService.On("SetStatusWallet", id, request).Return(&service.WalletResponse{
			WalletID:  1,
			Balance:   2000,
			Status:    "Deactive",
			CreatedAt: time.Date(2023, time.January, 27, 12, 30, 0, 0, time.UTC),
		}, nil)

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		r := `{"status":"Deactive"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(r))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/wallet/:id/status")
		c.SetParamNames("id")
		c.SetParamValues("1")

		expected := `{"wallet_id":1,"balance":2000,"status":"Deactive","created_at":"2023-01-27T12:30:00Z"}`

		// Assert
		if assert.NoError(t, walletHandler.ChangeStatus(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
		}
	})

	t.Run("change status id error", func(t *testing.T) {
		// Arrange
		request := service.StatusWalletRequest{
			Status: "Deactive",
		}
		walletService := service.NewWalletServiceMock()
		walletService.On("SetStatusWallet", "abc", request).Return(&service.WalletResponse{}, errs.NewBadRequest("id must be number"))

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		r := `{"status":"Deactive"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(r))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/wallet/:id/status")
		c.SetParamNames("id")
		c.SetParamValues("abc")

		// Assert
		if assert.NoError(t, walletHandler.ChangeStatus(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("change status request incorrect", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		request := service.StatusWalletRequest{
			Status: "Deactive",
		}
		walletService := service.NewWalletServiceMock()
		walletService.On("SetStatusWallet", id, request).Return(&service.WalletResponse{}, errs.NewBadRequest("request body incorrect format"))

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		r := `{"status":1}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(r))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/wallet/:id/status")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// Assert
		if assert.NoError(t, walletHandler.ChangeStatus(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("change status unexpected error", func(t *testing.T) {
		// Arrange
		var id int64 = 1
		request := service.StatusWalletRequest{
			Status: "Deactive",
		}
		walletService := service.NewWalletServiceMock()
		walletService.On("SetStatusWallet", id, request).Return(&service.WalletResponse{}, errs.NewUnexpectedError())

		walletHandler := handler.NewWalletHandler(walletService)

		// Act
		r := `{"status":"Deactive"}`
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(r))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/wallet/:id/status")
		c.SetParamNames("id")
		c.SetParamValues("1")

		// Assert
		if assert.NoError(t, walletHandler.ChangeStatus(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)
		}
	})
}
