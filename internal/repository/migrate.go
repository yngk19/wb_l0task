package repository

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/yngk19/wb_l0task/internal/config"
)

func Migrate(cfg config.DB) error {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode
	)
	m, err := migrate.New(sourceUrl, DBUrl)
	if err != nil {
		return err
	}
	m.Steps(2)
	return nil
}
