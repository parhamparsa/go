package persist

import (
	"app/internal/domain/entity"
	interfaces "app/internal/domain/interface"
	"app/pkg/sqlt"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

type userRepository struct {
	store *MysqlStore
	table *sqlt.SqlTable
}

var _ interfaces.UserRepository = userRepository{}

func newUserRepo(s *MysqlStore) interfaces.UserRepository {
	return &userRepository{
		store: s,
		table: sqlt.NewTable("users", "ID", []string{
			"id", "first_name", "last_name", "email", "active"},
		),
	}
}

func (r userRepository) Create(ctx context.Context, user *entity.User) (int64, error) {
	query, args, err := r.store.qb().
		Insert(r.table.Name()).
		Columns("first_name", "last_name", "email", "active").
		Values(user.FirstName, user.LastName, user.Email, user.Active).
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

func (r userRepository) find(ctx context.Context, id uint32, selectForUpdate bool) (*entity.User, error) {
	qb := r.store.qb().Select(r.table.Columns()...).From(r.table.Name()).Where(sq.Eq{"id": id})
	if selectForUpdate {
		qb = qb.Suffix("FOR UPDATE")
	}
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.store.masterDB.QueryRowContext(ctx, query, args...)

	var campaign entity.User
	err = row.Scan(userFields(&campaign)...)
	return &campaign, err
}

func (r userRepository) Find(ctx context.Context, id uint32) (*entity.User, error) {
	return r.find(ctx, id, false)
}

func (r userRepository) FindAndLock(ctx context.Context, id uint32) (*entity.User, error) {
	return r.find(ctx, id, true)
}
