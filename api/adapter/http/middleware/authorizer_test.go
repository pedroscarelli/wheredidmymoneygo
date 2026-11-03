package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
	dto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/account/request"
)

var secret = []byte("test-secret")

type fakeUseCase struct{ known map[string]bool }

func (f fakeUseCase) Create(*dto.CreateAccountRequestDTO) (*domain.PublicAccount, error) {
	return nil, nil
}
func (f fakeUseCase) Delete(string) error                               { return nil }
func (f fakeUseCase) Update(string, *dto.UpdateAccountRequestDTO) error { return nil }
func (f fakeUseCase) Movement(*dto.MovementRequestDTO) error            { return nil }
func (f fakeUseCase) GetByID(id string) (*domain.PublicAccount, error) {
	if !f.known[id] {
		return nil, domain.ErrAccountNotFound
	}
	return &domain.PublicAccount{ID: id}, nil
}

func sign(t *testing.T, method jwt.SigningMethod, claims jwt.MapClaims) string {
	token, err := jwt.NewWithClaims(method, claims).SignedString(secret)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func TestJwtAuthorizer(t *testing.T) {
	future := time.Now().Add(time.Hour).Unix()
	past := time.Now().Add(-time.Hour).Unix()

	tests := []struct {
		name       string
		header     string
		wantStatus int
	}{
		{"valid token", "Bearer " + sign(t, jwt.SigningMethodHS256, jwt.MapClaims{"sub": "acc-1", "exp": future}), http.StatusOK},
		{"missing header", "", http.StatusUnauthorized},
		{"not bearer", "Basic abc", http.StatusUnauthorized},
		{"expired", "Bearer " + sign(t, jwt.SigningMethodHS256, jwt.MapClaims{"sub": "acc-1", "exp": past}), http.StatusUnauthorized},
		{"no expiration", "Bearer " + sign(t, jwt.SigningMethodHS256, jwt.MapClaims{"sub": "acc-1"}), http.StatusUnauthorized},
		{"empty subject", "Bearer " + sign(t, jwt.SigningMethodHS256, jwt.MapClaims{"sub": "", "exp": future}), http.StatusUnauthorized},
		{"other algorithm", "Bearer " + sign(t, jwt.SigningMethodHS512, jwt.MapClaims{"sub": "acc-1", "exp": future}), http.StatusUnauthorized},
		{"deleted account", "Bearer " + sign(t, jwt.SigningMethodHS256, jwt.MapClaims{"sub": "gone", "exp": future}), http.StatusUnauthorized},
		{"garbage", "Bearer not.a.token", http.StatusUnauthorized},
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(UserClaimsKey).(*Claims)
		if claims.AccountID != "acc-1" {
			t.Errorf("claims account = %q", claims.AccountID)
		}
		w.WriteHeader(http.StatusOK)
	})
	handler := JwtAuthorizer(secret, fakeUseCase{known: map[string]bool{"acc-1": true}})(next)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/account/acc-1", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
