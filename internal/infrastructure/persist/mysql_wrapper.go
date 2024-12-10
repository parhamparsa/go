package persist

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"

	"github.com/Masterminds/squirrel"
)

const deadlockErrorCode = 1213

type DBWrapper struct {
	db                   *sql.DB
	maxRetriesOnDeadLock int
}

func newDbWrapper(db *sql.DB) *DBWrapper {
	return &DBWrapper{
		db:                   db,
		maxRetriesOnDeadLock: 3, //to be passed through config
	}
}

func (d *DBWrapper) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	return tx, err
}

func (d *DBWrapper) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	db := d.dbSelector(ctx)
	var result sql.Result
	var err error
	for i := 0; i < d.maxRetriesOnDeadLock; i++ {
		result, err = db.ExecContext(ctx, query, args...)
		if isDeadlockError(err) == false {
			break
		}
	}
	return result, err
}

func isDeadlockError(err error) bool {
	if err == nil {
		return false
	}
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == deadlockErrorCode
}

type DbExecutor interface {
	squirrel.ExecerContext
	squirrel.QueryerContext
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (d *DBWrapper) dbSelector(ctx context.Context) DbExecutor {
	value := ctx.Value(txKey)
	if txWrapper, ok := value.(*TxWrapper); ok && txWrapper.IsActive() {
		return txWrapper.tx
	}
	return d.db
}
func (d *DBWrapper) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	db := d.dbSelector(ctx)
	rows, err := db.QueryContext(ctx, query, args...)
	return rows, err
}
func (d *DBWrapper) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	db := d.dbSelector(ctx)
	return db.QueryRowContext(ctx, query, args...)
}
