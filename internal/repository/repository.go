package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	postgresql "github.com/yngk19/wb_l0task/internal/repository/client/postgres"
	"github.com/yngk19/wb_l0task/internal/repository/models"
	"log/slog"
)

type Repository struct {
	client postgresql.Client
	logger *slog.Logger
}

func (r *Repository) SaveOrder(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO orders (data)
		VALUES ($1)
		RETURNING id
	`
	row := r.client.QueryRow(ctx, query, order.Data)
	if err := row.Scan(&order.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState())
			r.logger.Error(newErr.Error())
			return newErr
		}
		return err
	}
	return nil
}

func (r *Repository) GetOrders(ctx context.Context) (u []models.Order, err error) {
	query := `
		SELECT id, data FROM orders
	`
	rows, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	orders := make([]models.Order, 0)
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.ID, &order.Data)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *Repository) GetOneOrder(ctx context.Context, id int) (models.Order, error) {
	query := `
		SELECT id, data 
		FROM orders 
		WHERE id = $1
	`
	var order models.Order
	row := r.client.QueryRow(ctx, query, id)
	if err := row.Scan(&order.ID, &order.Data); err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (r *Repository) GetOrdersByLimit(ctx context.Context, limit int) ([]models.Order, error) {
	query := `
		SELECT id, data
		FROM orders
		ORDER BY id
		DESC LIMIT $1
	`
	rows, err := r.client.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	orders := make([]models.Order, 0)
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.ID, &order.Data)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func NewRepository(client postgresql.Client, logger *slog.Logger) *Repository {
	return &Repository{
		client: client,
		logger: logger,
	}
}
