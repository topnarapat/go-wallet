package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/topnarapat/go-wallet/handler"
	"github.com/topnarapat/go-wallet/repository"
	"github.com/topnarapat/go-wallet/service"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("connect to database error", err)
	}
	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	walletRepositoryDB := repository.NewWalletRepository(db)
	walletService := service.NewWalletService(walletRepositoryDB)
	walletHandler := handler.NewWalletHandler(walletService)

	e.GET("/wallet", walletHandler.ListWallets)
	e.GET("/wallet/:id", walletHandler.GetWallet)
	e.POST("/wallet", walletHandler.CreateWallet)
	e.PUT("/wallet/:id", walletHandler.AddBalance)
	e.PUT("/wallet/:id/status", walletHandler.ChangeStatus)

	go func() {
		if err := e.Start(":2565"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
