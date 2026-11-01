package accountrepository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
	dto "github.com/pedroscarelli/wheredidmymoneygo/core/dto/account/request"
)

const uniqueViolation = "23505"

func (repository repository) Create(accountRequest *dto.CreateAccountRequestDTO, createdAt time.Time) (*domain.PublicAccount, error) {
	ctx := context.Background()
	account := domain.PublicAccount{}

	err := repository.db.QueryRow(
		ctx,
		`INSERT INTO tb_account
		(pk_st_id, st_first_name, st_last_name, st_cpf, st_email, st_password, dt_created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING pk_st_id, st_first_name, st_last_name, st_cpf, st_email, it_number, dt_created_at`,
		uuid.New().String(),
		accountRequest.FirstName,
		accountRequest.LastName,
		accountRequest.CPF,
		accountRequest.Email,
		accountRequest.Password,
		createdAt,
	).Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.CPF,
		&account.Email,
		&account.Number,
		&account.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return nil, domain.ErrDuplicateAccount
		}
		return nil, err
	}

	return &account, nil
}
