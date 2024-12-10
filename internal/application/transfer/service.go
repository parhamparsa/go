package transfer

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

var _ application.TransferServiceInterface = Service{}

func NewService(store interfaces.Store) Service {
	return Service{store: store}
}

func (s Service) Transfer(ctx context.Context, transfer *entity.Transfer) error {
	return s.store.ExecInTx(ctx, func(ctx context.Context) error {
		if transfer.To == transfer.From {
			return fmt.Errorf("sender and receiver cannot be the same")
		}

		fromAccount, err := s.store.Account().FindByUserIdAndLock(ctx, uint32(transfer.From))
		if err != nil {
			return err
		}
		if fromAccount.Balance < transfer.Amount {
			return fmt.Errorf("transfer amount is bigger that corrent value of user")
		}

		toAccount, err := s.store.Account().FindByUserIdAndLock(ctx, uint32(transfer.To))
		if toAccount == nil {
			return fmt.Errorf("transfer account not found")
		}

		newBalanceFrom := fromAccount.Balance - transfer.Amount
		newBalanceTo := toAccount.Balance + transfer.Amount

		err = s.store.Account().UpdateAccountBalance(ctx, uint32(fromAccount.UserID), newBalanceFrom)
		if err != nil {
			return err
		}

		err = s.store.Account().UpdateAccountBalance(ctx, uint32(toAccount.UserID), newBalanceTo)
		if err != nil {
			return err
		}

		err = s.store.Transfer().AddTransferRecord(ctx, transfer)
		if err != nil {
			return err
		}

		return nil
	})
}
