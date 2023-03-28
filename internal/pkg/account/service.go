package account

import (
	"context"

	"github.com/rudiarta/boilerplate-go/internal/pkg/model"
	"github.com/rudiarta/boilerplate-go/internal/pkg/repository"
)

type accountUsecaseCtx struct {
	accountRepo repository.AccountRepository
}

type AccountUsecase interface {
	CreateAccount(ctx context.Context, account model.Account) error
}

func NewAccountUsecase(
	accountRepo repository.AccountRepository,
) AccountUsecase {
	return &accountUsecaseCtx{
		accountRepo: accountRepo,
	}
}
