package domain

import "time"

type Customer struct {
	ID        string
	UniqueID  string
	Name      string
	Phone     string
	CreatedAt time.Time
}
