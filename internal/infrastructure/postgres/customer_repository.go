package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"kapoortrader-loyalty/internal/domain"
	"kapoortrader-loyalty/internal/repository"
)

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) repository.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) GetByPhone(ctx context.Context, phone string) (*domain.Customer, error) {
	query := `SELECT id, unique_id, name, phone, created_at FROM customers WHERE phone = $1`

	var c domain.Customer
	err := r.db.QueryRowContext(ctx, query, phone).Scan(&c.ID, &c.UniqueID, &c.Name, &c.Phone, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil for both if no customer is found
		}
		return nil, fmt.Errorf("get customer by phone: %w", err)
	}
	return &c, nil
}

func (r *customerRepository) Create(ctx context.Context, c *domain.Customer) error {
	query := `INSERT INTO customers (id, unique_id, name, phone, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, c.ID, c.UniqueID, c.Name, c.Phone, c.CreatedAt)
	if err != nil {
		return fmt.Errorf("create customer: %w", err)
	}
	return nil
}
