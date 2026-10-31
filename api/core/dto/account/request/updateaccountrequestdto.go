package dto

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type UpdateAccountRequestDTO struct {
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
	Password  *string `json:"password,omitempty"`
}

func (d *UpdateAccountRequestDTO) Validate() error {
	if d.FirstName == nil && d.LastName == nil && d.Password == nil {
		return errors.New("nothing to update")
	}
	if d.FirstName != nil {
		trimmed := strings.TrimSpace(*d.FirstName)
		d.FirstName = &trimmed
		if err := validateName("firstName", trimmed); err != nil {
			return err
		}
	}
	if d.LastName != nil {
		trimmed := strings.TrimSpace(*d.LastName)
		d.LastName = &trimmed
		if err := validateName("lastName", trimmed); err != nil {
			return err
		}
	}
	if d.Password != nil {
		return validatePassword(*d.Password)
	}

	return nil
}

func FromJSONUpdateAccountRequestDTO(body io.Reader) (*UpdateAccountRequestDTO, error) {
	updateAccountRequestDto := UpdateAccountRequestDTO{}
	if err := json.NewDecoder(body).Decode(&updateAccountRequestDto); err != nil {
		return nil, errors.New("invalid request body")
	}
	if err := updateAccountRequestDto.Validate(); err != nil {
		return nil, err
	}

	return &updateAccountRequestDto, nil
}
