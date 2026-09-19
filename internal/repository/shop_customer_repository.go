package repository

import (
	"context"

	"kapoortrader-loyalty/internal/domain"
)

type ShopCustomerRepository interface {
	Get(
		ctx context.Context,
		shopID string,
		customerID string,
	) (*domain.ShopCustomer, error)

	Create(ctx context.Context, shopCustomer *domain.ShopCustomer) error

	IncrementPurchase(
		ctx context.Context,
		shopID string,
		customerID string,
	) error

	ListByShop(
		ctx context.Context,
		shopID string,
	) ([]domain.ShopCustomer, error)

	ListDashboard(ctx context.Context, shopID string) ([]domain.DashboardRow, error)

	RedeemReward(ctx context.Context, shopID, customerID string) error
}
