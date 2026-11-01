package di

import (
	"github.com/pedroscarelli/wheredidmymoneygo/adapter/http/accountservice"
	"github.com/pedroscarelli/wheredidmymoneygo/adapter/postgres"
	"github.com/pedroscarelli/wheredidmymoneygo/adapter/postgres/accountrepository"
	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
	"github.com/pedroscarelli/wheredidmymoneygo/core/usecase/accountusecase"
)

func ConfigAccountDI(conn postgres.PoolInterface) (domain.AccountService, domain.AccountUseCase) {
	accountRepository := accountrepository.New(conn)
	accountUseCase := accountusecase.New(accountRepository)
	accountService := accountservice.New(accountUseCase)

	return accountService, accountUseCase
}
