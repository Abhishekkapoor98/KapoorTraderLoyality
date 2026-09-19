package domain

import "time"

type Shopkeeper struct {
	ID           string
	ShopID       string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}
