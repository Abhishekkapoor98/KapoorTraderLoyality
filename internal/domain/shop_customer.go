package domain

import "time"

type ShopCustomer struct {
	ShopID              string
	CustomerID          string
	PurchaseCount       int
	RewardCount         int
	RedeemedRewardCount int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (sc ShopCustomer) AvailableRewards() int {
	return sc.RewardCount - sc.RedeemedRewardCount
}
