package http

import (
	"encoding/json"
	"net/http"
	"text/template"

	"kapoortrader-loyalty/internal/usecase"
)

type CustomerHandler struct {
	customerService *usecase.CustomerService
}

func NewCustomerHandler(cs *usecase.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: cs}
}

type LoginRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type LoginResponse struct {
	CustomerID    string `json:"customer_id"`
	Name          string `json:"name"`
	PurchaseCount int    `json:"purchase_count"`
	Rewards       int    `json:"rewards"`
	Message       string `json:"message"`
}

func (h *CustomerHandler) LoginAndPurchase(w http.ResponseWriter, r *http.Request) {
	// Go 1.22+ supports extracting wildcards directly from the URL path
	shopID := r.PathValue("shopID")
	if shopID == "" {
		http.Error(w, "missing shop ID", http.StatusBadRequest)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Phone == "" {
		http.Error(w, "name and phone are required", http.StatusBadRequest)
		return
	}

	// Pass the data to our Clean Architecture use case
	shopCustomer, customer, err := h.customerService.LoginAndPurchase(r.Context(), shopID, req.Name, req.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := "Purchase recorded!"
	if shopCustomer.AvailableRewards() > 0 { // Business logic remains safely in the domain model
		msg = "🎁 Free Prize Available!"
	}

	resp := LoginResponse{
		CustomerID:    customer.UniqueID,
		Name:          customer.Name,
		PurchaseCount: shopCustomer.PurchaseCount,
		Rewards:       shopCustomer.AvailableRewards(),
		Message:       msg,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *CustomerHandler) ServeCustomerLoginPage(w http.ResponseWriter, r *http.Request) {
	shopID := r.PathValue("shopID")
	if shopID == "" {
		http.Error(w, "missing shop ID", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/customer_login.html")
	if err != nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
