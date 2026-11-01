package accountrepository

import (
	"github.com/pedroscarelli/wheredidmymoneygo/adapter/postgres"
	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
)

type repository struct {
	db postgres.PoolInterface
}

func New(db postgres.PoolInterface) domain.AccountRepository {
	return &repository{
		db: db,
	}
}
