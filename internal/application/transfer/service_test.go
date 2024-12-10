package transfer_test

import (
	"app/internal/domain/entity"
	"context"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestTransferService(t *testing.T) {
	tb.Reset(t)
	tb.Reset(t)
	err := tb.App().GetUserService().Create(context.Background(), &entity.User{
		ID: 1, FirstName: "John", LastName: "doe", Email: "johndoe@gmail.com", Active: true,
	})
	require.NoError(t, err)

	err = tb.App().GetUserService().Create(context.Background(), &entity.User{
		ID: 2, FirstName: "Mark", LastName: "Grossman", Email: "markgrossman@gmail.com", Active: true,
	})
	require.NoError(t, err)

	err = tb.App().GetAccountService().UpdateAccountBalance(context.Background(), 1, 1000)
	require.NoError(t, err)

	tests := map[string]struct {
		fromAccount                int32
		toAccount                  int32
		amounts                    []int64
		FromAccountExpectedBalance int64
		ToAccountExpectedBalance   int64
	}{
		"process_budget_correctly": {
			fromAccount:                1,
			toAccount:                  2,
			amounts:                    []int64{100, 50, 1, 2, 3, 6, 5, 2, 3, 4, 5, 7, 8, 9, 10},
			FromAccountExpectedBalance: 1000 - 215,
			ToAccountExpectedBalance:   215,
		},
	}

	for scenario, test := range tests {
		t.Run(scenario, func(t *testing.T) {
			var wg sync.WaitGroup
			for _, amountToTransfer := range test.amounts {
				wg.Add(1)
				go func() {
					defer wg.Done()
					err = tb.App().GetTransferService().
						Transfer(context.Background(), &entity.Transfer{
							From: test.fromAccount, To: test.toAccount, Amount: amountToTransfer})
					require.NoError(t, err)
				}()
			}
			wg.Wait()
			fromAccount, err := tb.App().GetAccountService().Find(context.Background(), uint32(test.fromAccount))
			require.NoError(t, err)
			require.Equal(t, test.FromAccountExpectedBalance, fromAccount.Balance)

			toAccount, err := tb.App().GetAccountService().Find(context.Background(), uint32(test.toAccount))
			require.NoError(t, err)
			require.Equal(t, test.ToAccountExpectedBalance, toAccount.Balance)

		})
	}
}
