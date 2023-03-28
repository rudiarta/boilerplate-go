package registry

import (
	"github.com/rudiarta/boilerplate-go/internal/app/config"
	"github.com/rudiarta/boilerplate-go/internal/pkg/repository"
)

type repositoryRegistry struct {
	cfg config.ConfigCtx
}

type RepositoryRegistry interface {
	AccountRepository() repository.AccountRepository
}

func NewRepositoryRegistry(cfg config.ConfigCtx) RepositoryRegistry {
	return &repositoryRegistry{
		cfg: cfg,
	}
}

func (c *repositoryRegistry) AccountRepository() repository.AccountRepository {
	return repository.NewAccountRepository(c.cfg)
}
