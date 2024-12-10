package account

import (
	"app/internal/application"
	"app/internal/domain/entity"
	interfaces "app/internal/domain/interface"
	"context"
	"fmt"
	"log"
)

type Service struct {
	store interfaces.Store
}

var _ application.AccountServiceInterface = Service{}

func NewService(store interfaces.Store) Service {
	return Service{store: store}
}

func (s Service) Create(ctx context.Context, account *entity.Account) (int64, error) {
	if err := account.Validate(); err != nil {
		return 0, fmt.Errorf("validation failed: %w", err)
	}

	id, err := s.store.Account().Create(ctx, account)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return id, nil
}

func (s Service) UpdateAccountBalance(ctx context.Context, accountID uint32, newBalance int64) error {
	if err := s.store.Account().UpdateAccountBalance(ctx, accountID, newBalance); err != nil {
		return err
	}
	return nil
}

func (s Service) Find(ctx context.Context, id uint32) (*entity.Account, error) {
	account, err := s.store.Account().Find(ctx, id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s Service) FindByUserId(ctx context.Context, userId uint32) (*entity.Account, error) {
	account, err := s.store.Account().FindByUserIdAndLock(ctx, userId)
	if err != nil {
		return nil, err
	}
	return account, nil
}
