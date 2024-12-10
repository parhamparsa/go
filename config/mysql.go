package config

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type MySql struct {
	Host                      string `yaml:"host"`
	User                      string `yaml:"user"`
	Password                  string `yaml:"password"`
	DBName                    string `yaml:"dbName"`
	MaxOpenConnections        int    `yaml:"maxOpenConnections"`
	MaxIdleConnections        int    `yaml:"maxIdleConnections"`
	MaxConnectionLifeTime     int    `yaml:"maxConnectionLifeTime"`
	MaxIdleConnectionLifeTime int    `yaml:"maxIdleConnectionLifeTime"`
}

func (c MySql) MigrationTable() string {
	return "migrations"
}
func (c MySql) LockTable() string {
	return "db_lock"
}

func (c MySql) MigrationLockKey() string {
	return fmt.Sprintf("migrations_lock")
}

func (c MySql) GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true&interpolateParams=true",
		c.User,
		c.Password,
		c.Host,
		c.DBName,
	)
}
func (c MySql) GetDSNWithoutDb() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/",
		c.User,
		c.Password,
		c.Host,
	)
}
func (c MySql) GetMigrationDSN() (string, error) {
	cfg, err := mysql.ParseDSN(c.GetDSN())
	if err != nil {
		panic(err)
	}
	if cfg.Params == nil {
		cfg.Params = map[string]string{}
	}

	cfg.Params["multiStatements"] = "true"
	cfg.ReadTimeout = 0

	return cfg.FormatDSN(), nil
}

func (c MySql) MaxConnectionLifeTimeSeconds() time.Duration {
	return time.Duration(c.MaxConnectionLifeTime) * time.Second
}

func (c MySql) MaxIdleConnectionLifeTimeSeconds() time.Duration {
	return time.Duration(c.MaxIdleConnectionLifeTime) * time.Second
}
