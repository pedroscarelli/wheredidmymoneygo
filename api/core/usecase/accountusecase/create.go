package accountusecase

import (
	"time"

	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
	dto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/account/request"
	"github.com/pedroscarelli/wheredidmymoneygo/core/security"
)

func (usecase usecase) Create(createAccountRequestDto *dto.CreateAccountRequestDTO) (*domain.PublicAccount, error) {
	hashedPassword, err := security.HashPassword(createAccountRequestDto.Password)
	if err != nil {
		return nil, err
	}

	createAccountRequestDto.Password = hashedPassword

	return usecase.repository.Create(createAccountRequestDto, time.Now().UTC())
}
