//go:build integration
// +build integration

package handler_test

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/topnarapat/go-wallet/handler"
	"github.com/topnarapat/go-wallet/repository"
	"github.com/topnarapat/go-wallet/service"
)

const serverPort = 2565

func TestGetAllWalletsIntegration(t *testing.T) {
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/wallets?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		walletRepo := repository.NewWalletRepository(db)
		walletSrv := service.NewWalletService(walletRepo)
		walletHandler := handler.NewWalletHandler(walletSrv)

		e.GET("/wallet", walletHandler.ListWallets)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := ``
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/wallet", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	expected := `[{"wallet_id":1,"balance":1000,"status":"Active","created_at":"2023-01-27T12:30:00Z"},{"wallet_id":2,"balance":2000,"status":"Active","created_at":"2023-01-27T12:30:00Z"},{"wallet_id":3,"balance":3000,"status":"Active","created_at":"2023-01-27T12:30:00Z"},{"wallet_id":4,"balance":4000,"status":"Active","created_at":"2023-01-27T12:30:00Z"},{"wallet_id":5,"balance":5000,"status":"Active","created_at":"2023-01-27T12:30:00Z"}]`

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestGetWalletIntegration(t *testing.T) {
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/wallets?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		walletRepo := repository.NewWalletRepository(db)
		walletSrv := service.NewWalletService(walletRepo)
		walletHandler := handler.NewWalletHandler(walletSrv)

		e.GET("/wallet/:id", walletHandler.GetWallet)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := ``
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/wallet/1", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	expected := `{"wallet_id":1,"balance":1000,"status":"Active","created_at":"2023-01-27T12:30:00Z"}`

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestCreateWalletIntegration(t *testing.T) {
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/wallets?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		walletRepo := repository.NewWalletRepository(db)
		walletSrv := service.NewWalletService(walletRepo)
		walletHandler := handler.NewWalletHandler(walletSrv)

		e.POST("/wallet", walletHandler.CreateWallet)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := `{"balance":6000}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/wallet", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	// Assertions

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestAddBalanceIntegration(t *testing.T) {
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/wallets?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		walletRepo := repository.NewWalletRepository(db)
		walletSrv := service.NewWalletService(walletRepo)
		walletHandler := handler.NewWalletHandler(walletSrv)

		e.PUT("/wallet/:id", walletHandler.AddBalance)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := `{"balance":1000,"operation":"Add"}`
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:%d/wallet/1", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	expected := `{"wallet_id":1,"balance":2000,"status":"Active","created_at":"2023-01-27T12:30:00Z"}`

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestChangeStatusIntegration(t *testing.T) {
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/wallets?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		walletRepo := repository.NewWalletRepository(db)
		walletSrv := service.NewWalletService(walletRepo)
		walletHandler := handler.NewWalletHandler(walletSrv)

		e.PUT("/wallet/:id/status", walletHandler.ChangeStatus)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	// Arrange
	reqBody := `{"status":"Deactive"}`
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:%d/wallet/2/status", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	expected := `{"wallet_id":2,"balance":2000,"status":"Deactive","created_at":"2023-01-27T12:30:00Z"}`

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}
