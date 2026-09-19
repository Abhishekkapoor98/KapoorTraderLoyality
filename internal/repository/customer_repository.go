package repository

import (
	"context"

	"kapoortrader-loyalty/internal/domain"
)

type CustomerRepository interface {
	GetByPhone(ctx context.Context, phone string) (*domain.Customer, error)
	Create(ctx context.Context, customer *domain.Customer) error
}
