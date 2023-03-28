package registry

import "github.com/rudiarta/boilerplate-go/internal/pkg/account"

type usecaseRegistry struct {
	repositoryRegistry RepositoryRegistry
}

type UsecaseRegistry interface {
	AccountUsecase() account.AccountUsecase
}

func NewUsecaseRegistry(
	repositoryRegistry RepositoryRegistry,
) UsecaseRegistry {
	return &usecaseRegistry{
		repositoryRegistry: repositoryRegistry,
	}
}

func (c *usecaseRegistry) AccountUsecase() account.AccountUsecase {
	return account.NewAccountUsecase(c.repositoryRegistry.AccountRepository())
}
