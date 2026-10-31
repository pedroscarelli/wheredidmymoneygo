package dto

import (
	"encoding/json"
	"errors"
	"io"
	"net/mail"
	"strings"
	"unicode"
)

const MinPasswordLength = 8

type CreateAccountRequestDTO struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	CPF       string `json:"cpf"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (d *CreateAccountRequestDTO) Validate() error {
	d.FirstName = strings.TrimSpace(d.FirstName)
	d.LastName = strings.TrimSpace(d.LastName)
	d.Email = strings.TrimSpace(d.Email)

	if err := validateName("firstName", d.FirstName); err != nil {
		return err
	}
	if err := validateName("lastName", d.LastName); err != nil {
		return err
	}
	if len(d.CPF) != 11 || strings.IndexFunc(d.CPF, func(r rune) bool { return !unicode.IsDigit(r) }) != -1 {
		return errors.New("cpf must contain exactly 11 digits")
	}
	if _, err := mail.ParseAddress(d.Email); err != nil || len(d.Email) > 100 {
		return errors.New("email is invalid")
	}
	return validatePassword(d.Password)
}

func validateName(field, value string) error {
	if value == "" || len(value) > 50 {
		return errors.New(field + " must have between 1 and 50 characters")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return errors.New("password must have at least 8 characters")
	}
	return nil
}

func FromJSONCreateAccountRequestDTO(body io.Reader) (*CreateAccountRequestDTO, error) {
	createAccountRequestDTO := CreateAccountRequestDTO{}
	if err := json.NewDecoder(body).Decode(&createAccountRequestDTO); err != nil {
		return nil, errors.New("invalid request body")
	}
	if err := createAccountRequestDTO.Validate(); err != nil {
		return nil, err
	}

	return &createAccountRequestDTO, nil
}
