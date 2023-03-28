package account

import (
	"context"
	"fmt"

	"github.com/rudiarta/boilerplate-go/internal/pkg/model"
	"github.com/rudiarta/boilerplate-go/pkg/libxendit"
)

func (c *accountUsecaseCtx) CreateAccount(ctx context.Context, account model.Account) error {
	libxendit.CreateVa()
	c.accountRepo.Insert(ctx, account)
	fmt.Println("CreateAccount: ", account.Name)
	return nil
}
