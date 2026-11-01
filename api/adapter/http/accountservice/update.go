package accountservice

import (
	"errors"
	"net/http"

	"github.com/pedroscarelli/wheredidmymoneygo/adapter/http/response"
	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
	dto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/account/request"
)

// Update always changes the authenticated account; the id comes from the token, never from the body.
func (service service) Update(w http.ResponseWriter, r *http.Request) {
	accountID, err := GetAccountIDFromToken(r)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	updateAccountRequestDto, err := dto.FromJSONUpdateAccountRequestDTO(r.Body)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := service.usecase.Update(accountID, updateAccountRequestDto); err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		response.InternalError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "account updated successfully"})
}
