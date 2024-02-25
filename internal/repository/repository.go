package repository

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/yngk19/wb_l0task/internal/config"
	"github.com/yngk19/wb_l0task/internal/repository/models"
	"log/slog"
)

type Repository struct {
	DB *sqlx.DB
}

func NewRepository(cfg config.DB, log *slog.Logger) (*Repository, error) {
	db, err := sqlx.Open(
		"postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.DBName, cfg.Password, cfg.SSLMode,
		),
	)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Repository{DB: db}, nil
}

func (r *Repository) SaveOrder(order models.Order) (int64, error) {
	jsonData, err := json.Marshal(order)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	res, err := r.DB.NamedExec(`INSERT INTO orders (data) VALUES (:data)`, map[string]interface{}{
		"data": jsonData,
	})
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	id, err := res.LastInsertId()
	return id, nil
}

func (r *Repository) GetOrderById(id int) (*models.Order, error) {
	return nil, nil
}
