package accountusecase

import (
	"errors"
	"testing"
	"time"

	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
	"github.com/pedroscarelli/wheredidmymoneygo/core/domain/enum"
	dto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/account/request"
	"golang.org/x/crypto/bcrypt"
)

type fakeRepository struct {
	accounts      map[string]*domain.Account
	updated       *domain.Account
	movementCents int
	movementErr   error
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{accounts: map[string]*domain.Account{
		"acc-1": {ID: "acc-1", FirstName: "Pedro", LastName: "Scarelli", Password: "old-hash"},
	}}
}

func (f *fakeRepository) Create(d *dto.CreateAccountRequestDTO, createdAt time.Time) (*domain.PublicAccount, error) {
	return &domain.PublicAccount{ID: "new", FirstName: d.FirstName}, nil
}
func (f *fakeRepository) Delete(accountID string) error { return nil }
func (f *fakeRepository) Update(account *domain.Account) error {
	copied := *account
	f.updated = &copied
	return nil
}
func (f *fakeRepository) GetAccountByID(accountID string) (*domain.Account, error) {
	account, ok := f.accounts[accountID]
	if !ok {
		return nil, domain.ErrAccountNotFound
	}
	copied := *account
	return &copied, nil
}
func (f *fakeRepository) GetAccountByEmail(email string) (*domain.Account, error) {
	return nil, domain.ErrAccountNotFound
}
func (f *fakeRepository) Movement(d *dto.MovementRequestDTO, valueInCents int, createdAt time.Time) error {
	f.movementCents = valueInCents
	return f.movementErr
}
func (f *fakeRepository) GetAccountBalance(accountID string) (int, error) { return 0, nil }

func ptr(s string) *string { return &s }

func TestUpdateKeepsFieldsThatWereNotSent(t *testing.T) {
	repo := newFakeRepository()
	uc := New(repo)

	if err := uc.Update("acc-1", &dto.UpdateAccountRequestDTO{FirstName: ptr("Pedro Henrique")}); err != nil {
		t.Fatal(err)
	}

	if repo.updated.FirstName != "Pedro Henrique" {
		t.Errorf("first name = %q, want %q", repo.updated.FirstName, "Pedro Henrique")
	}
	if repo.updated.LastName != "Scarelli" || repo.updated.Password != "old-hash" {
		t.Errorf("fields not sent were changed: %+v", repo.updated)
	}
}

func TestUpdateHashesNewPassword(t *testing.T) {
	t.Setenv("HASH_COST", "4")
	repo := newFakeRepository()

	if err := New(repo).Update("acc-1", &dto.UpdateAccountRequestDTO{Password: ptr("new-password")}); err != nil {
		t.Fatal(err)
	}

	if bcrypt.CompareHashAndPassword([]byte(repo.updated.Password), []byte("new-password")) != nil {
		t.Error("password was not stored as a bcrypt hash of the new password")
	}
}

func TestUpdateUnknownAccount(t *testing.T) {
	err := New(newFakeRepository()).Update("missing", &dto.UpdateAccountRequestDTO{FirstName: ptr("x")})
	if !errors.Is(err, domain.ErrAccountNotFound) {
		t.Fatalf("err = %v, want ErrAccountNotFound", err)
	}
}

func TestMovementSignAndRounding(t *testing.T) {
	tests := []struct {
		name  string
		kind  enum.MovementType
		value float64
		want  int
	}{
		{"credit", enum.CREDITO, 0.29, 29},
		{"debit is negative", enum.DEBITO, 19.99, -1999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepository()
			if err := New(repo).Movement(&dto.MovementRequestDTO{Type: tt.kind, Value: tt.value, Recurrence: 1}); err != nil {
				t.Fatal(err)
			}
			if repo.movementCents != tt.want {
				t.Errorf("cents = %d, want %d", repo.movementCents, tt.want)
			}
		})
	}
}
