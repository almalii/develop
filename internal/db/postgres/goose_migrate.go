package postgres

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/config"
	"notes-rew/internal/db/postgres/migrations"
)

func UpMigrations(c config.Config) error {
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		c.DB.UserName,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.DBName,
		c.DB.SSLMode,
	)

	gooseDB, err := goose.OpenDBWithDriver(c.DB.Driver, connStr)
	if err != nil {
		logrus.Error("error opening db connection on migrations")
		return err
	}

	defer gooseDB.Close()

	goose.SetBaseFS(&migrations.Content)

	err = goose.SetDialect(c.DB.Driver)
	if err != nil {
		logrus.Error("setting dialect error on migrations")
		return err
	}

	err = goose.Up(gooseDB, ".")
	if err != nil {
		logrus.Error("running migrations error on up")
		return err
	}

	return nil
}
