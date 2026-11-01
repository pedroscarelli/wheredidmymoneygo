package domain

import (
	"net/http"

	requestDto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/authentication/request"
	responseDto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/authentication/response"
)

type AuthenticationService interface {
	Login(response http.ResponseWriter, request *http.Request)
}

type AuthenticationUseCase interface {
	Login(accountRequest *requestDto.LoginRequestDTO) (*responseDto.LoginResponseDTO, error)
}
