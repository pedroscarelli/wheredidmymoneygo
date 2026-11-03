package dto

import (
	"testing"
	"time"

	enum "github.com/pedroscarelli/wheredidmymoneygo/core/domain/enum"
)

func validMovement() MovementRequestDTO {
	return MovementRequestDTO{
		Value:       10.5,
		Type:        enum.CREDITO,
		Recurrence:  1,
		DueDate:     time.Date(2025, 5, 27, 0, 0, 0, 0, time.UTC),
		Description: "salary",
	}
}

func TestMovementValidate(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*MovementRequestDTO)
		wantErr bool
	}{
		{"valid credit", func(d *MovementRequestDTO) {}, false},
		{"valid debit", func(d *MovementRequestDTO) { d.Type = enum.DEBITO }, false},
		{"invalid type", func(d *MovementRequestDTO) { d.Type = "PIX" }, true},
		{"zero value", func(d *MovementRequestDTO) { d.Value = 0 }, true},
		{"negative value", func(d *MovementRequestDTO) { d.Value = -5 }, true},
		{"value too big", func(d *MovementRequestDTO) { d.Value = 10_000_000.01 }, true},
		{"recurrence 0", func(d *MovementRequestDTO) { d.Recurrence = 0 }, true},
		{"recurrence 13", func(d *MovementRequestDTO) { d.Recurrence = 13 }, true},
		{"recurrence 12", func(d *MovementRequestDTO) { d.Recurrence = 12 }, false},
		{"missing due date", func(d *MovementRequestDTO) { d.DueDate = time.Time{} }, true},
		{"blank description", func(d *MovementRequestDTO) { d.Description = "   " }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := validMovement()
			tt.modify(&d)
			if err := d.Validate(); (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValueInCentsRounds(t *testing.T) {
	tests := []struct {
		value float64
		want  int
	}{
		{0.29, 29},
		{19.99, 1999},
		{0.1 + 0.2, 30},
		{100, 10000},
	}

	for _, tt := range tests {
		d := MovementRequestDTO{Value: tt.value}
		if got := d.ValueInCents(); got != tt.want {
			t.Errorf("ValueInCents(%v) = %d, want %d", tt.value, got, tt.want)
		}
	}
}
