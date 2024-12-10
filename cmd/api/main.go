package main

import (
	"app/config"
	"app/internal/application"
	"app/internal/application/account"
	"app/internal/application/transfer"
	"app/internal/application/user"
	"app/internal/infrastructure/persist"
	"app/internal/infrastructure/rest/v1"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	app, err := createCoreApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	restApi := v1.NewRestApiV1(cfg, app)
	log.Fatal(restApi.Run())
}

func createCoreApp(cfg *config.Config) (application.App, error) {
	store, err := persist.NewMysqlStore(&cfg.MySql)

	if err != nil {
		return nil, fmt.Errorf("hello unable to connect to mysql: %s", err.Error())
	}

	userService := user.NewService(store)
	accountService := account.NewService(store)
	transferService := transfer.NewService(store)

	return application.NewAppCore(
		userService,
		accountService,
		transferService,
	), nil
}
