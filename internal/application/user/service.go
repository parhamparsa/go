package user

import (
	"app/internal/application"
	"app/internal/domain/entity"
	interfaces "app/internal/domain/interface"
	"context"
	"fmt"
)

type Service struct {
	store interfaces.Store
}

var _ application.UserServiceInterface = Service{}

func NewService(store interfaces.Store) Service {
	return Service{store: store}
}

func (s Service) Create(ctx context.Context, user *entity.User) error {
	if err := user.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return s.store.ExecInTx(ctx, func(ctx context.Context) error {
		id, err := s.store.User().Create(ctx, user)
		if err != nil {
			return err
		}

		account := entity.Account{
			UserID:  id,
			Type:    entity.AccountTypeNormal,
			Balance: 0,
		}
		_, errAccount := s.store.Account().Create(ctx, &account)
		if errAccount != nil {
			return errAccount
		}
		return nil
	})
}

func (s Service) Find(ctx context.Context, id uint32) (*entity.User, error) {
	result, err := s.store.User().Find(ctx, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
