package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	enum "github.com/pedroscarelli/wheredidmymoneygo/core/domain/enum"
)

// maxValue keeps the value in cents inside a Postgres INTEGER column.
const maxValue = 10_000_000.00

type MovementRequestDTO struct {
	AccountID   string            `json:"-"`
	Value       float64           `json:"value"`
	Type        enum.MovementType `json:"type"`
	Recurrence  int               `json:"recurrence"`
	DueDate     time.Time         `json:"dueDate"`
	Description string            `json:"description"`
}

func (d *MovementRequestDTO) Validate() error {
	switch d.Type {
	case enum.DEBITO, enum.CREDITO:
	default:
		return fmt.Errorf("invalid type: %s. Allowed values: %s or %s",
			d.Type,
			enum.DEBITO,
			enum.CREDITO)
	}

	if math.IsNaN(d.Value) || d.Value <= 0 || d.Value > maxValue {
		return errors.New("value must be greater than 0 and at most 10000000")
	}

	if d.Recurrence < 1 || d.Recurrence > 12 {
		return errors.New("recurrence must be between 1 and 12")
	}

	if d.DueDate.IsZero() {
		return errors.New("dueDate is required")
	}

	d.Description = strings.TrimSpace(d.Description)
	if d.Description == "" || len(d.Description) > 255 {
		return errors.New("description must have between 1 and 255 characters")
	}

	return nil
}

// ValueInCents converts the value to cents, rounding instead of truncating
// (0.29 * 100 is 28.999... in floating point).
func (d *MovementRequestDTO) ValueInCents() int {
	return int(math.Round(d.Value * 100))
}

func FromJSONCreateMovementRequestDTO(body io.Reader) (*MovementRequestDTO, error) {
	createMovementRequestDTO := MovementRequestDTO{}
	if err := json.NewDecoder(body).Decode(&createMovementRequestDTO); err != nil {
		return nil, errors.New("invalid request body")
	}

	if err := createMovementRequestDTO.Validate(); err != nil {
		return nil, err
	}

	return &createMovementRequestDTO, nil
}
