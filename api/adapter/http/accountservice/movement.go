package accountservice

import (
	"net/http"

	"github.com/pedroscarelli/wheredidmymoneygo/adapter/http/response"
	dto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/account/request"
)

func (service service) Movement(w http.ResponseWriter, r *http.Request) {
	accountID, err := GetAccountIDFromToken(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	movementRequestDTO, err := dto.FromJSONCreateMovementRequestDTO(r.Body)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	movementRequestDTO.AccountID = accountID
	if err := service.usecase.Movement(movementRequestDTO); err != nil {
		response.InternalError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, map[string]string{"message": "movement registered successfully"})
}
