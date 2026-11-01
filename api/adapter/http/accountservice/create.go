package accountservice

import (
	"errors"
	"net/http"

	"github.com/pedroscarelli/wheredidmymoneygo/adapter/http/response"
	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
	dto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/account/request"
)

func (service service) Create(w http.ResponseWriter, r *http.Request) {
	accountRequest, err := dto.FromJSONCreateAccountRequestDTO(r.Body)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	account, err := service.usecase.Create(accountRequest)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateAccount) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		response.InternalError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, account)
}
