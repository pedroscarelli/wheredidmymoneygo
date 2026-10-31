package dto

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type LoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func FromJSONCreateLoginRequestDTO(body io.Reader) (*LoginRequestDTO, error) {
	loginRequestDTO := LoginRequestDTO{}
	if err := json.NewDecoder(body).Decode(&loginRequestDTO); err != nil {
		return nil, errors.New("invalid request body")
	}

	loginRequestDTO.Email = strings.TrimSpace(loginRequestDTO.Email)
	if loginRequestDTO.Email == "" || loginRequestDTO.Password == "" {
		return nil, errors.New("email and password are required")
	}

	return &loginRequestDTO, nil
}
