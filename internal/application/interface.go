package application

import (
	"app/internal/domain/entity"
	"context"
)

//go:generate mockgen -source=interface.go -destination=../../mocks/app.go
type App interface {
	GetUserService() UserServiceInterface
	GetAccountService() AccountServiceInterface
	GetTransferService() TransferServiceInterface
}

type UserServiceInterface interface {
	Create(ctx context.Context, account *entity.User) error
	Find(ctx context.Context, id uint32) (*entity.User, error)
}

type AccountServiceInterface interface {
	Create(ctx context.Context, account *entity.Account) (int64, error)
	UpdateAccountBalance(ctx context.Context, accountID uint32, newBalance int64) error
	Find(ctx context.Context, id uint32) (*entity.Account, error)
	FindByUserId(ctx context.Context, userId uint32) (*entity.Account, error)
}

type TransferServiceInterface interface {
	Transfer(ctx context.Context, transfer *entity.Transfer) error
}
