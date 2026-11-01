package authenticationservice

import (
	"errors"
	"net/http"

	"github.com/pedroscarelli/wheredidmymoneygo/adapter/http/response"
	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
	dto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/authentication/request"
)

func (service service) Login(w http.ResponseWriter, r *http.Request) {
	loginRequestDto, err := dto.FromJSONCreateLoginRequestDTO(r.Body)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	loginResponseDto, err := service.usecase.Login(loginRequestDto)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		response.InternalError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, loginResponseDto)
}
