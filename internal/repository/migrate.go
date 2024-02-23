package repository

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/yngk19/wb_l0task/internal/config"
	"log/slog"
)

func Migrate(cfg config.DB, log *slog.Logger) error {
	log.Info("Making migrations")
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)
	sourcURL := fmt.Sprintf("file:///%s", cfg.MigrationsPath)
	m, err := migrate.New(sourcURL, dbURL)
	if err != nil {
		return err
	}
	m.Steps(1)
	return nil

}
