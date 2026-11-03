package authenticationusecase

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
	accountDto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/account/request"
	dto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/authentication/request"
	"golang.org/x/crypto/bcrypt"
)

type fakeRepository struct {
	account *domain.Account
	err     error
}

func (f fakeRepository) Create(*accountDto.CreateAccountRequestDTO, time.Time) (*domain.PublicAccount, error) {
	return nil, nil
}
func (f fakeRepository) Delete(string) error                                           { return nil }
func (f fakeRepository) Update(*domain.Account) error                                  { return nil }
func (f fakeRepository) GetAccountByID(string) (*domain.Account, error)                { return f.account, f.err }
func (f fakeRepository) GetAccountByEmail(string) (*domain.Account, error)             { return f.account, f.err }
func (f fakeRepository) Movement(*accountDto.MovementRequestDTO, int, time.Time) error { return nil }
func (f fakeRepository) GetAccountBalance(string) (int, error)                         { return 0, nil }

func hash(t *testing.T, password string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	return string(h)
}

func TestLogin(t *testing.T) {
	t.Setenv("JWT_SECRET_KEY", "test-secret")
	account := &domain.Account{ID: "acc-1", Password: hash(t, "right-password")}

	tests := []struct {
		name     string
		repo     fakeRepository
		password string
		wantErr  error
	}{
		{"success", fakeRepository{account: account}, "right-password", nil},
		{"wrong password", fakeRepository{account: account}, "wrong-password", domain.ErrInvalidCredentials},
		{"unknown email", fakeRepository{err: domain.ErrAccountNotFound}, "right-password", domain.ErrInvalidCredentials},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := New(tt.repo).Login(&dto.LoginRequestDTO{Email: "a@b.com", Password: tt.password})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}

			claims := jwt.MapClaims{}
			if _, err := jwt.ParseWithClaims(res.Token, claims, func(*jwt.Token) (any, error) { return []byte("test-secret"), nil }); err != nil {
				t.Fatalf("token is not valid: %v", err)
			}
			if claims["sub"] != "acc-1" {
				t.Errorf("sub = %v, want acc-1", claims["sub"])
			}
		})
	}
}

func TestLoginWithoutSecretFails(t *testing.T) {
	t.Setenv("JWT_SECRET_KEY", "")
	account := &domain.Account{ID: "acc-1", Password: hash(t, "pw")}

	if _, err := New(fakeRepository{account: account}).Login(&dto.LoginRequestDTO{Email: "a@b.com", Password: "pw"}); err == nil {
		t.Fatal("expected an error when JWT_SECRET_KEY is empty")
	}
}
