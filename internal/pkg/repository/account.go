package repository

import (
	"context"
	"fmt"

	"github.com/rudiarta/boilerplate-go/internal/app/config"
	"github.com/rudiarta/boilerplate-go/internal/pkg/model"
)

type accountRepositoryCtx struct {
	cfg config.ConfigCtx
}

type AccountRepository interface {
	Insert(ctx context.Context, account model.Account) error
}

func NewAccountRepository(cfg config.ConfigCtx) AccountRepository {
	return &accountRepositoryCtx{
		cfg: cfg,
	}
}

func (c *accountRepositoryCtx) Insert(ctx context.Context, account model.Account) error {
	fmt.Println("Insert: ", account.Name)
	db, err := c.cfg.GormDbConnection()
	if err != nil {
	}
	d, _ := db.DB()
	err = d.Ping()
	if err != nil {
	}
	return nil
}
