package domain

type DashboardRow struct {
	CustomerID    string `json:"customer_id"`
	Name          string `json:"name"`
	UniqueID      string `json:"unique_id"`
	Phone         string `json:"phone"`
	PurchaseCount int    `json:"purchase_count"`
	Rewards       int    `json:"rewards"`
}
