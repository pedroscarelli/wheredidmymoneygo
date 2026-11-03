package accountservice

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
	dto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/account/request"
)

type fakeUseCase struct {
	createErr     error
	updatedID     string
	deletedID     string
	getErr        error
	movementCalls int
}

func (f *fakeUseCase) Create(*dto.CreateAccountRequestDTO) (*domain.PublicAccount, error) {
	return &domain.PublicAccount{ID: "new"}, f.createErr
}
func (f *fakeUseCase) Delete(id string) error { f.deletedID = id; return nil }
func (f *fakeUseCase) Update(id string, _ *dto.UpdateAccountRequestDTO) error {
	f.updatedID = id
	return nil
}
func (f *fakeUseCase) GetByID(id string) (*domain.PublicAccount, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return &domain.PublicAccount{ID: id}, nil
}
func (f *fakeUseCase) Movement(*dto.MovementRequestDTO) error { f.movementCalls++; return nil }

// serve routes the request like main.go does, with the given account already authenticated.
func serve(uc *fakeUseCase, method, path, body, accountID string) *httptest.ResponseRecorder {
	svc := New(uc)
	router := mux.NewRouter()
	router.HandleFunc("/account", svc.Create).Methods(http.MethodPost)
	router.HandleFunc("/account", svc.Update).Methods(http.MethodPatch)
	router.HandleFunc("/account/movement", svc.Movement).Methods(http.MethodPost)
	router.HandleFunc("/account/{id}", svc.GetByID).Methods(http.MethodGet)
	router.HandleFunc("/account/{id}", svc.Delete).Methods(http.MethodDelete)

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if accountID != "" {
		ctx := req.Context()
		req = req.WithContext(contextWithAccount(ctx, accountID))
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestGetByIDOnlyAllowsOwnAccount(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{"own account", "/account/acc-1", http.StatusOK},
		{"someone else's account", "/account/acc-2", http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(&fakeUseCase{}, http.MethodGet, tt.path, "", "acc-1")
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestDeleteOnlyAllowsOwnAccount(t *testing.T) {
	uc := &fakeUseCase{}
	if rec := serve(uc, http.MethodDelete, "/account/acc-2", "", "acc-1"); rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", rec.Code)
	}
	if uc.deletedID != "" {
		t.Fatalf("account %q was deleted by another user", uc.deletedID)
	}

	if rec := serve(uc, http.MethodDelete, "/account/acc-1", "", "acc-1"); rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestUpdateUsesAccountFromToken(t *testing.T) {
	uc := &fakeUseCase{}
	rec := serve(uc, http.MethodPatch, "/account", `{"id":"acc-2","firstName":"Hacker"}`, "acc-1")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", rec.Code, rec.Body)
	}
	if uc.updatedID != "acc-1" {
		t.Fatalf("updated account %q, want acc-1", uc.updatedID)
	}
}

func TestCreateStatusCodes(t *testing.T) {
	valid := `{"firstName":"Pedro","lastName":"Scarelli","cpf":"12345678901","email":"p@example.com","password":"secret123"}`

	tests := []struct {
		name       string
		body       string
		createErr  error
		wantStatus int
	}{
		{"created", valid, nil, http.StatusCreated},
		{"invalid json", "{", nil, http.StatusBadRequest},
		{"invalid fields", `{"firstName":"Pedro"}`, nil, http.StatusBadRequest},
		{"duplicate", valid, domain.ErrDuplicateAccount, http.StatusConflict},
		{"database error", valid, errors.New("pq: connection refused"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(&fakeUseCase{createErr: tt.createErr}, http.MethodPost, "/account", tt.body, "")
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if strings.Contains(rec.Body.String(), "connection refused") {
				t.Fatal("internal error details leaked to the client")
			}
		})
	}
}

func TestMovementRejectsInvalidBody(t *testing.T) {
	uc := &fakeUseCase{}
	rec := serve(uc, http.MethodPost, "/account/movement", `{"value":-1,"type":"CREDITO"}`, "acc-1")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	if uc.movementCalls != 0 {
		t.Fatal("invalid movement reached the use case")
	}
}

func TestProtectedRouteWithoutClaims(t *testing.T) {
	if rec := serve(&fakeUseCase{}, http.MethodGet, "/account/acc-1", "", ""); rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}
