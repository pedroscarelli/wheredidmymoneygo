package accountusecase

import (
	"time"

	"github.com/pedroscarelli/wheredidmymoneygo/core/domain/enum"
	dto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/account/request"
)

func (usecase usecase) Movement(movementRequestDto *dto.MovementRequestDTO) error {
	valueInCents := movementRequestDto.ValueInCents()
	if movementRequestDto.Type == enum.DEBITO {
		valueInCents = -valueInCents
	}

	return usecase.repository.Movement(movementRequestDto, valueInCents, time.Now().UTC())
}
