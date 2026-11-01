package authenticationusecase

import (
	"errors"

	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
	requestDto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/authentication/request"
	responseDto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/authentication/response"
	"github.com/pedroscarelli/wheredidmymoneygo/core/security"
)

func (usecase usecase) Login(loginRequestDto *requestDto.LoginRequestDTO) (*responseDto.LoginResponseDTO, error) {
	account, err := usecase.repository.GetAccountByEmail(loginRequestDto.Email)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if security.IsPasswordIncorrect(loginRequestDto.Password, account.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	jwtToken, err := security.GenerateJwtToken(account.ID)
	if err != nil {
		return nil, err
	}

	return &responseDto.LoginResponseDTO{Token: jwtToken}, nil
}
