package persist

import (
	"app/config"
	migrations "app/db"
	interfaces "app/internal/domain/interface"
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/mattermost/morph"
	"github.com/mattermost/morph/drivers/mysql"
	"github.com/mattermost/morph/sources/embedded"
	"github.com/pkg/errors"
	"path/filepath"
	"time"
)

type mysqlStoreRepos struct {
	userRepo     interfaces.UserRepository
	accountRepo  interfaces.AccountRepository
	transferRepo interfaces.TransferRepository
}

type MysqlStore struct {
	masterDB     *DBWrapper
	repositories *mysqlStoreRepos
	masterConfig *config.MySql
}

func NewMysqlStore(cfg *config.MySql) (interfaces.Store, error) {
	masterDB, err := newConn(cfg.GetDSN(), cfg)
	if err != nil {
		return nil, fmt.Errorf("host = %s, user = %s err=%w", cfg.Host, cfg.User, err)
	}

	store := &MysqlStore{
		masterDB:     newDbWrapper(masterDB),
		masterConfig: cfg,
	}

	// migrations will be run automatically
	err = store.migrate(cfg)
	if err != nil {
		return nil, fmt.Errorf("migrate err=%w", err)
	}

	store.repositories = &mysqlStoreRepos{
		userRepo:     newUserRepo(store),
		accountRepo:  newAccountRepo(store),
		transferRepo: newTransferRepo(store),
	}

	return store, nil
}

func (s MysqlStore) migrate(cfg *config.MySql) error {
	basedir := migrations.MigrationsBaseDir()

	fs := migrations.MigrationsFS()
	migrationNames, err := migrations.EntryNames(fs, basedir)
	if err != nil {
		return err
	}

	assetfunc := func(name string) ([]byte, error) {
		return fs.ReadFile(filepath.Join(basedir, name))
	}

	src, err := embedded.WithInstance(embedded.Resource(migrationNames, assetfunc))
	if err != nil {
		return err
	}

	dsn, err := cfg.GetMigrationDSN()
	if err != nil {
		return err
	}

	db, err := newConn(dsn, cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := mysql.WithInstance(db)
	if err != nil {
		return err
	}

	lockOption := morph.WithLock(cfg.MigrationLockKey())
	ctx := context.Background()
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer cancelFunc()
	engine, err := morph.New(
		timeout, driver, src, lockOption, morph.SetMigrationTableName(cfg.MigrationTable()))
	if err != nil {
		return fmt.Errorf(
			"could not initiate migration. Maybe locking failed or server is not alive. original err: %w", err)
	}
	defer engine.Close()
	return engine.ApplyAll()
}

func newConn(dsn string, cfg *config.MySql) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.MaxConnectionLifeTimeSeconds())
	db.SetConnMaxIdleTime(cfg.MaxIdleConnectionLifeTimeSeconds())
	return db, nil
}

func (s MysqlStore) qb() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Question)
}

func (s MysqlStore) ExecInTx(ctx context.Context, f func(ctx context.Context) error) error {
	ctx, tx, err := s.BeginTx(ctx)
	if err != nil {
		return err
	}
	if err = f(ctx); err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("error while rolling back another error %v", err)
		}
		return err
	}
	return tx.Commit()
}

func (s MysqlStore) BeginTx(ctx context.Context) (context.Context, driver.Tx, error) {
	value := ctx.Value(txKey)
	if txWrapper, ok := value.(TxWrapper); !ok || txWrapper.tx == nil {
		tx, err := s.masterDB.BeginTx(ctx)
		if err != nil {
			return ctx, nil, err
		}
		txWrapper := &TxWrapper{tx: tx}
		return context.WithValue(ctx, txKey, txWrapper), txWrapper, nil
	}
	if txWrapper, ok := value.(*TxWrapper); ok {
		return ctx, txWrapper, nil
	}
	return ctx, nil, errors.New("unexpected err in BeginTx")
}

func (t *TxWrapper) Commit() error {
	if !t.IsActive() {
		return errors.New("calling commit on closed tx")
	}
	err := t.tx.Commit()
	t.tx = nil
	return err
}

func (t *TxWrapper) Rollback() error {
	if !t.IsActive() {
		return fmt.Errorf("calling rollback on closed tx")
	}
	err := t.tx.Rollback()
	t.tx = nil
	return err
}

func (t *TxWrapper) IsActive() bool {
	return t.tx != nil
}

const txKey = "_tx"

type TxWrapper struct {
	tx *sql.Tx
}

// CreateDb todo limit this only to testing env or making it configurable
func CreateDb(cfg *config.MySql) error {
	db, err := newConn(cfg.GetDSNWithoutDb(), cfg)

	if err != nil {
		return err
	}
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cfg.DBName)
	timeout, cancelFunc := context.WithTimeoutCause(
		context.Background(), time.Second*2, errors.New("unable to create db"+cfg.DBName),
	)
	defer cancelFunc()
	row := db.QueryRowContext(timeout, query)
	return row.Err()
}

// TruncateAllTables todo limit this only to testing env or making it configurable
func (s MysqlStore) TruncateAllTables() (err error) {
	rows, err := s.masterDB.QueryContext(context.Background(), `show tables`)
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		closeErr := rows.Close()
		if err != nil || closeErr != nil {
			err = fmt.Errorf("err: %w:err on closing : %w", err, closeErr)
		}
	}(rows)
	for rows.Next() {
		var tableName string
		if err = rows.Scan(&tableName); err != nil {
			return err
		}
		if err = s.TruncateTable(tableName); err != nil {
			return err
		}
	}

	return rows.Err()
}

func (s MysqlStore) TruncateTable(tableName string) error {
	if tableName != s.masterConfig.MigrationTable() && tableName != s.masterConfig.LockTable() {
		query := fmt.Sprintf("TRUNCATE TABLE %s", tableName)
		_, err := s.masterDB.ExecContext(context.Background(), query)
		return err
	}

	return nil
}

func (r MysqlStore) User() interfaces.UserRepository {
	return r.repositories.userRepo
}

func (r MysqlStore) Account() interfaces.AccountRepository {
	return r.repositories.accountRepo
}

func (r MysqlStore) Transfer() interfaces.TransferRepository {
	return r.repositories.transferRepo
}
