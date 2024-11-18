package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/guregu/null/v5"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"

	"basic-go-project/pkg/clients"
	"basic-go-project/pkg/clients/postgresql"
	"basic-go-project/src/entities/customErrors"
	"basic-go-project/src/entities/dtos"
)

type accountPgStorage struct {
	client postgresql.Client
	logger *zerolog.Logger
}

func NewAccountStoragePG(client postgresql.Client, logger *zerolog.Logger) clients.AccountStorage {
	return &accountPgStorage{
		client: client,
		logger: logger,
	}
}

type AccountDB struct {
	ID       null.String `json:"id"`
	Login    null.String `json:"login"`
	Password null.String `json:"password"`
	IsActive null.Bool   `json:"isActive"`
}

func (a accountPgStorage) Create(ctx context.Context, account *dtos.Account) error {
	query := `
		INSERT INTO service.accounts
		    (
		     login, 
		     password,
		     is_active
		     )
		VALUES
		       ($1, $2, $3)
		RETURNING id
	`

	if err := a.client.QueryRow(ctx, query,
		account.Login,
		account.Password,
		account.IsActive,
	).Scan(&account.ID); err != nil {
		return fmt.Errorf("client.QueryRow() failed: %w", err)
	}

	return nil
}

func (a accountPgStorage) FindAll(ctx context.Context, limit, offset int) ([]dtos.Account, error) {
	query := `
		SELECT id, 
		       login, 
		       password, 
		       is_active 
		FROM service.accounts
		ORDER BY id
		LIMIT $1 OFFSET $2;
	`

	rows, err := a.client.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("client.Query() failed: %w", err)
	}
	defer rows.Close()

	accountsList := make([]dtos.Account, 0)
	for rows.Next() {
		var accountDB AccountDB

		err = rows.Scan(
			&accountDB.ID,
			&accountDB.Login,
			&accountDB.Password,
			&accountDB.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("client.Query() failed: %w", err)
		}

		account := dtos.Account{
			ID:       accountDB.ID.String,
			Login:    accountDB.Login.String,
			Password: accountDB.Password.String,
			IsActive: accountDB.IsActive.Bool,
		}

		accountsList = append(accountsList, account)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("client.Query() failed: %w", err)
	}

	return accountsList, nil
}

var ErrNoAccountFound = errors.New("didn't find account")

func (a accountPgStorage) FindOne(ctx context.Context, id string) (dtos.Account, error) {
	query := `
		SELECT id, 
		       login, 
		       password, 
		       is_active
		FROM service.accounts
		WHERE id = $1;
	`

	var accountDB AccountDB
	err := a.client.QueryRow(ctx, query, id).Scan(
		&accountDB.ID,
		&accountDB.Login,
		&accountDB.Password,
		&accountDB.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dtos.Account{}, ErrNoAccountFound
		}

		return dtos.Account{}, fmt.Errorf("client.QueryRow() failed: %w", err)
	}

	account := dtos.Account{
		ID:       accountDB.ID.String,
		Login:    accountDB.Login.String,
		Password: accountDB.Password.String,
		IsActive: accountDB.IsActive.Bool,
	}

	return account, nil
}

func (a accountPgStorage) Update(ctx context.Context, account dtos.Account) error {
	query := `
		   UPDATE service.accounts 
		   SET  
		       login = $1, 
		       password = $2, 
		       is_active = $3
           WHERE id = $4;
`

	_, err := a.client.Exec(ctx, query,
		account.Login,
		account.Password,
		account.IsActive,
		account.ID,
	)
	if err != nil {
		return fmt.Errorf("client.Query() failed: %w", err)
	}

	return nil
}

func (a accountPgStorage) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM service.accounts WHERE id=$1;
	`

	commandTag, err := a.client.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("client.Exec() failed: %w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return customErrors.ErrNoRowsFindToDelete
	}

	a.logger.Debug().Msgf("account with id = %s sucsefuly DELETED", id)

	return nil
}
