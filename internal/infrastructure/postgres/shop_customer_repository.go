package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"kapoortrader-loyalty/internal/domain"
	"kapoortrader-loyalty/internal/repository"
)

type shopCustomerRepository struct {
	db *sql.DB
}

func NewShopCustomerRepository(db *sql.DB) repository.ShopCustomerRepository {
	return &shopCustomerRepository{db: db}
}

func (r *shopCustomerRepository) Get(ctx context.Context, shopID, customerID string) (*domain.ShopCustomer, error) {
	query := `SELECT shop_id, customer_id, purchase_count, reward_count, redeemed_reward_count, created_at, updated_at 
              FROM shop_customers WHERE shop_id = $1 AND customer_id = $2`

	var sc domain.ShopCustomer
	err := r.db.QueryRowContext(ctx, query, shopID, customerID).Scan(
		&sc.ShopID, &sc.CustomerID, &sc.PurchaseCount, &sc.RewardCount,
		&sc.RedeemedRewardCount, &sc.CreatedAt, &sc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get shop customer: %w", err)
	}
	return &sc, nil
}

func (r *shopCustomerRepository) Create(ctx context.Context, sc *domain.ShopCustomer) error {
	query := `INSERT INTO shop_customers (shop_id, customer_id, purchase_count, reward_count, redeemed_reward_count, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, sc.ShopID, sc.CustomerID, sc.PurchaseCount, sc.RewardCount, sc.RedeemedRewardCount, sc.CreatedAt, sc.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create shop customer: %w", err)
	}
	return nil
}

func (r *shopCustomerRepository) IncrementPurchase(ctx context.Context, shopID, customerID string) error {
	// Automatically calculates the new reward_count in the database (1 reward per 5 purchases)
	query := `UPDATE shop_customers 
              SET purchase_count = purchase_count + 1, 
                  reward_count = (purchase_count + 1) / 5, 
                  updated_at = CURRENT_TIMESTAMP 
              WHERE shop_id = $1 AND customer_id = $2`
	_, err := r.db.ExecContext(ctx, query, shopID, customerID)
	if err != nil {
		return fmt.Errorf("increment purchase: %w", err)
	}
	return nil
}

func (r *shopCustomerRepository) ListByShop(ctx context.Context, shopID string) ([]domain.ShopCustomer, error) {
	return nil, nil // We will implement this later when building the Shopkeeper Dashboard
}

func (r *shopCustomerRepository) ListDashboard(ctx context.Context, shopID string) ([]domain.DashboardRow, error) {
	query := `
        SELECT c.id, c.name, c.unique_id, c.phone, sc.purchase_count, (sc.reward_count - sc.redeemed_reward_count) as rewards
        FROM shop_customers sc
        JOIN customers c ON sc.customer_id = c.id
        WHERE sc.shop_id = $1
        ORDER BY sc.updated_at DESC
    `
	rows, err := r.db.QueryContext(ctx, query, shopID)
	if err != nil {
		return nil, fmt.Errorf("list dashboard: %w", err)
	}
	defer rows.Close()

	var dashboard []domain.DashboardRow
	for rows.Next() {
		var row domain.DashboardRow
		// Fixed: Now we return the error if scanning fails, and append if it succeeds.
		if err := rows.Scan(&row.CustomerID, &row.Name, &row.UniqueID, &row.Phone, &row.PurchaseCount, &row.Rewards); err != nil {
			return nil, fmt.Errorf("scan dashboard row: %w", err)
		}
		dashboard = append(dashboard, row)
	}

	// Fixed: Moved outside the loop
	if dashboard == nil {
		dashboard = []domain.DashboardRow{}
	}

	return dashboard, nil
}

// Added the missing RedeemReward method needed for the UI
func (r *shopCustomerRepository) RedeemReward(ctx context.Context, shopID, customerID string) error {
	query := `
        UPDATE shop_customers 
        SET redeemed_reward_count = redeemed_reward_count + 1, 
            updated_at = CURRENT_TIMESTAMP 
        WHERE shop_id = $1 AND customer_id = $2 AND reward_count > redeemed_reward_count
    `
	res, err := r.db.ExecContext(ctx, query, shopID, customerID)
	if err != nil {
		return fmt.Errorf("redeem reward: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no rewards available to redeem")
	}
	return nil
}
