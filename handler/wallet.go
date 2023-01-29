package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/topnarapat/go-wallet/errs"
	"github.com/topnarapat/go-wallet/service"
)

type walletHandler struct {
	walletSrv service.WalletService
}

type Err struct {
	Message string `json:"message"`
}

func NewWalletHandler(walletSrv service.WalletService) walletHandler {
	return walletHandler{walletSrv: walletSrv}
}

func handlerError(c echo.Context, err error) error {
	switch e := err.(type) {
	case errs.AppError:
		return c.JSON(e.Code, Err{Message: e.Message})
	default:
		return c.JSON(http.StatusInternalServerError, e)
	}
}

func (h walletHandler) ListWallets(c echo.Context) error {
	wallets, err := h.walletSrv.ListAllWallets()
	if err != nil {
		return handlerError(c, err)
	}

	return c.JSON(http.StatusOK, wallets)
}

func (h walletHandler) GetWallet(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handlerError(c, errs.NewBadRequest("id must be number"))
	}

	wallet, err := h.walletSrv.GetWalletDetail(int64(id))
	if err != nil {
		return handlerError(c, err)
	}

	return c.JSON(http.StatusOK, wallet)
}

func (h walletHandler) CreateWallet(c echo.Context) error {
	balance := service.WalletRequest{}
	err := c.Bind(&balance)
	if err != nil {
		return handlerError(c, errs.NewBadRequest("request body incorrect format"))
	}

	wallet, err := h.walletSrv.CreateWallet(balance)
	if err != nil {
		return handlerError(c, err)
	}

	return c.JSON(http.StatusCreated, wallet)
}

func (h walletHandler) AddBalance(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handlerError(c, errs.NewBadRequest("id must be number"))
	}

	amount := service.AddWalletRequest{}
	err = c.Bind(&amount)
	if err != nil {
		return handlerError(c, errs.NewBadRequest("request body incorrect format"))
	}

	wallet, err := h.walletSrv.SetWalletBalance(int64(id), amount)
	if err != nil {
		return handlerError(c, err)
	}

	return c.JSON(http.StatusOK, wallet)
}

func (h walletHandler) ChangeStatus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handlerError(c, errs.NewBadRequest("id must be number"))
	}

	status := service.StatusWalletRequest{}
	err = c.Bind(&status)
	if err != nil {
		return handlerError(c, errs.NewBadRequest("request body incorrect format"))
	}

	wallet, err := h.walletSrv.SetStatusWallet(int64(id), status)
	if err != nil {
		return handlerError(c, err)
	}

	return c.JSON(http.StatusOK, wallet)
}
