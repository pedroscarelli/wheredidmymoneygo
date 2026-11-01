package accountrepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
)

func (repository repository) GetAccountByEmail(email string) (*domain.Account, error) {
	ctx := context.Background()
	row := repository.db.QueryRow(ctx,
		`select 
		pk_st_id, st_first_name, st_last_name, st_cpf, st_email, it_number, dt_created_at, st_password
		from tb_account
		where st_email = $1`,
		email)

	account, err := scanIntoAccount(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrAccountNotFound
	}

	return account, err
}
