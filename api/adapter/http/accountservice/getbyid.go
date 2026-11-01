package accountservice

import (
	"errors"
	"net/http"

	"github.com/pedroscarelli/wheredidmymoneygo/adapter/http/response"
	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
)

func (service service) GetByID(w http.ResponseWriter, r *http.Request) {
	accountID, status := authorizedAccountID(r)
	if status != http.StatusOK {
		response.Error(w, status, http.StatusText(status))
		return
	}

	account, err := service.usecase.GetByID(accountID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			response.Error(w, http.StatusNotFound, err.Error())
			return
		}
		response.InternalError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, account)
}
