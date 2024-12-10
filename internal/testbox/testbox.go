package testbox

import (
	"app/config"
	"app/internal/application"
	"app/internal/application/account"
	"app/internal/application/transfer"
	"app/internal/application/user"
	interfaces "app/internal/domain/interface"
	"app/internal/infrastructure/persist"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestBox struct {
	cfg *config.Config
	s   interfaces.Store
	app application.App
}

var (
	tb *TestBox //nolint:gochecknoglobals
)

func SetupTestBox(packageName string) *TestBox {
	var err error
	tb, err = instance(packageName)
	if err != nil {
		log.Fatal(err)
	}
	err = tb.Store().TruncateAllTables() //nolint:all
	if err != nil {
		log.Fatal("Failed to truncate all tables", err)
	}
	return tb
}

func instance(packageName string) (*TestBox, error) {
	var (
		err   error
		cfg   *config.Config
		store interfaces.Store
	)

	cfg, err = makeTestConfiguration(packageName)
	if err != nil {
		return nil, err
	}
	err = persist.CreateDb(&cfg.MySql)
	if err != nil {
		return nil, err
	}

	store, err = persist.NewMysqlStore(&cfg.MySql)
	if err != nil {
		return nil, err
	}

	tb = &TestBox{
		cfg: cfg,
		s:   store,
		app: createApp(store),
	}

	return tb, err
}

func createApp(store interfaces.Store) application.App {
	userService := user.NewService(store)
	accountService := account.NewService(store)
	transferService := transfer.NewService(store)

	return application.NewAppCore(
		userService,
		accountService,
		transferService,
	)
}

func makeTestConfiguration(packageName string) (*config.Config, error) {
	cfg, err := config.NewConfig()
	cfg.MySql.DBName = packageName
	return cfg, err
}

func (tb *TestBox) Config() *config.Config {
	return tb.cfg
}

func (tb *TestBox) Store() interfaces.Store {
	return tb.s
}

func (tb *TestBox) App() application.App {
	return tb.app
}

func (tb *TestBox) Reset(t *testing.T) *TestBox {
	require.NoError(t, tb.Store().TruncateAllTables()) //nolint:all
	return tb
}
