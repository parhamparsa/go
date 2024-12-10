package persist

import (
	"app/internal/domain/entity"
	interfaces "app/internal/domain/interface"
	"app/pkg/sqlt"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

type accountRepository struct {
	store *MysqlStore
	table *sqlt.SqlTable
}

var _ interfaces.AccountRepository = accountRepository{}

func newAccountRepo(s *MysqlStore) interfaces.AccountRepository {
	return &accountRepository{
		store: s,
		table: sqlt.NewTable("accounts", "ID", []string{
			"id", "user_id", "type", "balance"},
		),
	}
}

func (r accountRepository) Create(ctx context.Context, account *entity.Account) (int64, error) {
	query, args, err := r.store.qb().
		Insert(r.table.Name()).
		Columns("user_id", "type", "balance").
		Values(account.UserID, account.Type, account.Balance).
		ToSql()

	if err != nil {
		return 0, err
	}

	result, err := r.store.masterDB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	generatedID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last inserted ID: %w", err)
	}

	return generatedID, nil
}

func (r accountRepository) UpdateAccountBalance(ctx context.Context, accountID uint32, newBalance int64) error {
	query, args, err := r.store.qb().
		Update(r.table.Name()).
		Set("balance", newBalance).
		Where(sq.Eq{"user_id": accountID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.store.masterDB.ExecContext(ctx, query, args...)
	return err
}

func (r accountRepository) find(
	ctx context.Context, id uint32, selectForUpdate bool) (*entity.Account, error) {
	qb := r.store.qb().Select(r.table.Columns()...).From(r.table.Name()).Where(sq.Eq{"id": id})
	if selectForUpdate {
		qb = qb.Suffix("FOR UPDATE")
	}
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.store.masterDB.QueryRowContext(ctx, query, args...)

	var account entity.Account
	err = row.Scan(accountFields(&account)...)
	return &account, err
}

func (r accountRepository) Find(ctx context.Context, id uint32) (*entity.Account, error) {
	return r.find(ctx, id, false)
}

func (r accountRepository) FindAndLock(ctx context.Context, id uint32) (*entity.Account, error) {
	return r.find(ctx, id, true)
}

func (r accountRepository) FindByUserIdAndLock(ctx context.Context, userId uint32) (*entity.Account, error) {
	qb := r.store.qb().Select(r.table.Columns()...).
		From(r.table.Name()).Where(sq.Eq{"user_id": userId}).Suffix(" FOR UPDATE")
	query, args, err := qb.ToSql()
	row := r.store.masterDB.QueryRowContext(ctx, query, args...)

	var account entity.Account
	err = row.Scan(accountFields(&account)...)
	return &account, err
}
