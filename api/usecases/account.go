package usecases

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"basic-go-project/pkg/clients"
	"basic-go-project/src/entities/dtos"
)

type AccountUseCase struct {
	db  clients.AccountStorage
	log *zerolog.Logger
}

func NewAccountUseCase(storage clients.AccountStorage, l *zerolog.Logger) AccountUseCase {
	return AccountUseCase{
		db:  storage,
		log: l,
	}
}

func (a AccountUseCase) Create(account *dtos.Account) error {
	ctx := context.Background()
	if err := a.db.Create(ctx, account); err != nil {
		return fmt.Errorf("db.Create(): %w", err)
	}

	return nil
}

func (a AccountUseCase) GetAll(limit, offset int) ([]dtos.Account, error) {
	ctx := context.Background()
	accountsList, err := a.db.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("db.FindAll(): %w", err)
	}

	return accountsList, nil
}

func (a AccountUseCase) GetByID(id string) (dtos.Account, error) {
	ctx := context.Background()
	account, err := a.db.FindOne(ctx, id)
	if err != nil {
		return dtos.Account{}, fmt.Errorf("a.db.FindOne(): %w", err)
	}

	return account, nil
}

func (a AccountUseCase) Update(updateAccountData *dtos.UpdateAccountRequest) error {
	ctx := context.Background()

	account, err := a.db.FindOne(ctx, updateAccountData.ID)
	if err != nil {
		return fmt.Errorf("a.db.FindOne(): %w", err)
	}

	if updateAccountData.Login != nil {
		account.Login = *updateAccountData.Login
	}
	if updateAccountData.Password != nil {
		account.Password = *updateAccountData.Password
	}
	if updateAccountData.IsActive != nil {
		account.IsActive = *updateAccountData.IsActive
	}

	if err := a.db.Update(ctx, account); err != nil {
		return fmt.Errorf("db.Update(): %w", err)
	}

	return nil
}
func (a AccountUseCase) Delete(id string) error {
	ctx := context.Background()
	if err := a.db.Delete(ctx, id); err != nil {
		return fmt.Errorf("db.Delete(): %w", err)
	}

	return nil
}
