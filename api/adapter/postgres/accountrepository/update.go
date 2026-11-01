package accountrepository

import (
	"context"
	"fmt"

	"github.com/pedroscarelli/wheredidmymoneygo/core/domain"
)

func (repository repository) Update(account *domain.Account) error {
	ctx := context.Background()

	query := `
		UPDATE tb_account
		SET st_first_name = $1,
		    st_last_name  = $2,
		    st_password   = $3
		WHERE pk_st_id = $4;
	    `

	result, err := repository.db.Exec(
		ctx,
		query,
		account.FirstName,
		account.LastName,
		account.Password,
		account.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating account in the database: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrAccountNotFound
	}

	return nil
}
