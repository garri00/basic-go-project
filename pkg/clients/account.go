package clients

import (
	"context"

	"basic-go-project/src/entities/dtos"
)

// TODO: add orm
type AccountStorage interface {
	Create(ctx context.Context, account *dtos.Account) error
	FindAll(ctx context.Context, limit, offset int) ([]dtos.Account, error)
	FindOne(ctx context.Context, id string) (dtos.Account, error)
	Update(ctx context.Context, account dtos.Account) error
	Delete(ctx context.Context, id string) error
}
