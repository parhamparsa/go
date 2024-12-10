package interfaces

import (
	"app/internal/domain/entity"
	"context"
	"database/sql"
)

type SQLStore struct {
	db *sql.DB
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

type Store interface {
	ExecInTx(ctx context.Context, f func(ctx context.Context) error) error
	TruncateAllTables() (err error)

	User() UserRepository
	Account() AccountRepository
	Transfer() TransferRepository
}

type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) (int64, error)
	UpdateAccountBalance(ctx context.Context, accountID uint32, newBalance int64) error
	Find(ctx context.Context, id uint32) (*entity.Account, error)
	FindAndLock(ctx context.Context, id uint32) (*entity.Account, error)
	FindByUserIdAndLock(ctx context.Context, userId uint32) (*entity.Account, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (int64, error)
	Find(ctx context.Context, id uint32) (*entity.User, error)
	FindAndLock(ctx context.Context, id uint32) (*entity.User, error)
}

type TransferRepository interface {
	AddTransferRecord(ctx context.Context, transfer *entity.Transfer) error
}
