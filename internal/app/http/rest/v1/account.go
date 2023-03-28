package v1

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rudiarta/boilerplate-go/internal/app/config"
	"github.com/rudiarta/boilerplate-go/internal/pkg/account"
	"github.com/rudiarta/boilerplate-go/internal/pkg/model"
)

type accountHandler struct {
	cfg            config.ConfigCtx
	accountUsecase account.AccountUsecase
}

func NewAccountHandler(
	cfg config.ConfigCtx,
	accountUsecase account.AccountUsecase,
) accountHandler {
	return accountHandler{
		cfg:            cfg,
		accountUsecase: accountUsecase,
	}
}

func (c *accountHandler) CreateAccount(e echo.Context) error {
	data := account.CreateAccountReq{
		Name: "test",
	}
	c.accountUsecase.CreateAccount(context.TODO(), model.Account{
		Name: data.Name,
	})
	return e.String(http.StatusOK, "Hello, World!")
}
