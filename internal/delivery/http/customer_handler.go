package http

import (
	"encoding/json"
	"html/template"
	"net/http"

	loyalty "kapoortrader-loyalty"
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

	shopCustomer, customer, err := h.customerService.LoginAndPurchase(r.Context(), shopID, req.Name, req.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := "Purchase recorded!"
	if shopCustomer.AvailableRewards() > 0 {
		msg = "🎉 Free Prize Available!"
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

	htmlBytes, err := loyalty.TemplatesFS.ReadFile("web/templates/customer_login.html")
	if err != nil {
		http.Error(w, "failed to read template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("login").Parse(string(htmlBytes))
	if err != nil {
		http.Error(w, "failed to parse template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "failed to execute template: "+err.Error(), http.StatusInternalServerError)
	}
}
