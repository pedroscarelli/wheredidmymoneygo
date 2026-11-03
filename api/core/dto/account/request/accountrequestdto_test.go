package dto

import (
	"strings"
	"testing"
)

func TestCreateAccountValidate(t *testing.T) {
	valid := func() CreateAccountRequestDTO {
		return CreateAccountRequestDTO{FirstName: "Pedro", LastName: "Scarelli", CPF: "12345678901", Email: "pedro@example.com", Password: "secret123"}
	}

	tests := []struct {
		name    string
		modify  func(*CreateAccountRequestDTO)
		wantErr bool
	}{
		{"valid", func(d *CreateAccountRequestDTO) {}, false},
		{"missing first name", func(d *CreateAccountRequestDTO) { d.FirstName = " " }, true},
		{"long last name", func(d *CreateAccountRequestDTO) { d.LastName = strings.Repeat("a", 51) }, true},
		{"short cpf", func(d *CreateAccountRequestDTO) { d.CPF = "123" }, true},
		{"cpf with letters", func(d *CreateAccountRequestDTO) { d.CPF = "1234567890a" }, true},
		{"invalid email", func(d *CreateAccountRequestDTO) { d.Email = "not-an-email" }, true},
		{"short password", func(d *CreateAccountRequestDTO) { d.Password = "1234567" }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := valid()
			tt.modify(&d)
			if err := d.Validate(); (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateAccountValidate(t *testing.T) {
	name := "Pedro"
	blank := " "
	short := "123"

	tests := []struct {
		name    string
		dto     UpdateAccountRequestDTO
		wantErr bool
	}{
		{"only first name", UpdateAccountRequestDTO{FirstName: &name}, false},
		{"empty body", UpdateAccountRequestDTO{}, true},
		{"blank last name", UpdateAccountRequestDTO{LastName: &blank}, true},
		{"short password", UpdateAccountRequestDTO{Password: &short}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.dto.Validate(); (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFromJSONRejectsInvalidBody(t *testing.T) {
	if _, err := FromJSONCreateAccountRequestDTO(strings.NewReader("{not json")); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
