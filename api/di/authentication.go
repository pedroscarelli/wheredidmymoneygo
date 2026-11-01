package di

import (
	"github.com/pedroscarelli/wheredidmymoneygo/adapter/http/authenticationservice"
	"github.com/pedroscarelli/wheredidmymoneygo/adapter/postgres"
	"github.com/pedroscarelli/wheredidmymoneygo/adapter/postgres/accountrepository"
	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
	"github.com/pedroscarelli/wheredidmymoneygo/core/usecase/authenticationusecase"
)

func ConfigAuthenticationDI(conn postgres.PoolInterface) domain.AuthenticationService {
	accountRepository := accountrepository.New(conn)
	authenticationUseCase := authenticationusecase.New(accountRepository)
	authenticationService := authenticationservice.New(authenticationUseCase)

	return authenticationService
}
