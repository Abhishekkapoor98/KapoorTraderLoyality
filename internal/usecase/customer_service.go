package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"kapoortrader-loyalty/internal/domain"
	"kapoortrader-loyalty/internal/repository"

	"github.com/google/uuid"
)

type CustomerService struct {
	customerRepo     repository.CustomerRepository
	shopCustomerRepo repository.ShopCustomerRepository
}

func NewCustomerService(cr repository.CustomerRepository, scr repository.ShopCustomerRepository) *CustomerService {
	return &CustomerService{
		customerRepo:     cr,
		shopCustomerRepo: scr,
	}
}

// LoginAndPurchase executes the primary customer workflow
func (s *CustomerService) LoginAndPurchase(ctx context.Context, shopID, name, phone string) (*domain.ShopCustomer, *domain.Customer, error) {
	// 1. Check if the customer already exists by phone[cite: 1]
	customer, err := s.customerRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, nil, err
	}

	// 2. If new customer, generate unique ID and create[cite: 1]
	if customer == nil {
		shortID := strings.ToUpper(strings.Split(uuid.New().String(), "-")[0])

		customer = &domain.Customer{
			ID:        uuid.New().String(),
			UniqueID:  fmt.Sprintf("CUS-%s", shortID),
			Name:      name,
			Phone:     phone,
			CreatedAt: time.Now(),
		}
		if err := s.customerRepo.Create(ctx, customer); err != nil {
			return nil, nil, err
		}
	}

	// 3. Find the customer's relationship with this specific shop[cite: 1]
	shopCustomer, err := s.shopCustomerRepo.Get(ctx, shopID, customer.ID)
	if err != nil {
		return nil, nil, err
	}

	// 4. Create shop relationship if it's their first time at this shop, otherwise increment[cite: 1]
	if shopCustomer == nil {
		shopCustomer = &domain.ShopCustomer{
			ShopID:              shopID,
			CustomerID:          customer.ID,
			PurchaseCount:       1,
			RewardCount:         0,
			RedeemedRewardCount: 0,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}
		if err := s.shopCustomerRepo.Create(ctx, shopCustomer); err != nil {
			return nil, nil, err
		}
	} else {
		if err := s.shopCustomerRepo.IncrementPurchase(ctx, shopID, customer.ID); err != nil {
			return nil, nil, err
		}
		// Increment the local object to return updated data
		shopCustomer.PurchaseCount++
		shopCustomer.RewardCount = shopCustomer.PurchaseCount / 5
	}

	return shopCustomer, customer, nil
}
