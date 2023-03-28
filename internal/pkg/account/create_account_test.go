package account

import (
	"context"
	"testing"

	"github.com/rudiarta/boilerplate-go/internal/pkg/model"
	"github.com/rudiarta/boilerplate-go/internal/pkg/repository"
	"github.com/rudiarta/boilerplate-go/internal/pkg/repository/mocks"
)

func Test_accountUsecaseCtx_CreateAccount(t *testing.T) {
	type fields struct {
		accountRepo repository.AccountRepository
	}
	type args struct {
		ctx     context.Context
		account model.Account
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			fields: func() fields {
				accountRepo := new(mocks.AccountRepository)
				accountRepo.On("Insert", context.TODO(), model.Account{}).Return(nil).Once()
				return fields{
					accountRepo: accountRepo,
				}
			}(),
			args: args{
				ctx: context.TODO(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &accountUsecaseCtx{
				accountRepo: tt.fields.accountRepo,
			}
			if err := c.CreateAccount(tt.args.ctx, tt.args.account); (err != nil) != tt.wantErr {
				t.Errorf("accountUsecaseCtx.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
