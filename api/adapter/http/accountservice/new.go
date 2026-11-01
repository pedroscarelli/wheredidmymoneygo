package accountservice

import "github.com/pedroscarelli/wheredidmymoneygo/core/domain"

type service struct {
	usecase domain.AccountUseCase
}

func New(usecase domain.AccountUseCase) domain.AccountService {
	return &service{
		usecase: usecase,
	}
}
