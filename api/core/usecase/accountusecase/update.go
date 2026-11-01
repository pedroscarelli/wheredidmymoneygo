package accountusecase

import (
	dto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/account/request"
	"github.com/pedroscarelli/wheredidmymoneygo/core/security"
)

// Update applies only the fields that were sent, keeping the rest of the account as it is.
func (usecase usecase) Update(accountID string, updateAccountRequestDto *dto.UpdateAccountRequestDTO) error {
	account, err := usecase.repository.GetAccountByID(accountID)
	if err != nil {
		return err
	}

	if updateAccountRequestDto.FirstName != nil {
		account.FirstName = *updateAccountRequestDto.FirstName
	}

	if updateAccountRequestDto.LastName != nil {
		account.LastName = *updateAccountRequestDto.LastName
	}

	if updateAccountRequestDto.Password != nil {
		hashedPassword, err := security.HashPassword(*updateAccountRequestDto.Password)
		if err != nil {
			return err
		}
		account.Password = hashedPassword
	}

	return usecase.repository.Update(account)
}
