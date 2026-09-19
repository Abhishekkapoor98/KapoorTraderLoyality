package http

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"kapoortrader-loyalty/internal/repository"

	"github.com/skip2/go-qrcode"
)

type ShopkeeperHandler struct {
	shopCustomerRepo repository.ShopCustomerRepository
}

func NewShopkeeperHandler(scr repository.ShopCustomerRepository) *ShopkeeperHandler {
	return &ShopkeeperHandler{shopCustomerRepo: scr}
}

func (h *ShopkeeperHandler) GetDashboardData(w http.ResponseWriter, r *http.Request) {
	shopID := r.PathValue("shopID")
	if shopID == "" {
		http.Error(w, "missing shop ID", http.StatusBadRequest)
		return
	}

	dashboard, err := h.shopCustomerRepo.ListDashboard(r.Context(), shopID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboard)
}

func (h *ShopkeeperHandler) ServeDashboardPage(w http.ResponseWriter, r *http.Request) {
	shopID := r.PathValue("shopID")
	if shopID == "" {
		http.Error(w, "missing shop ID", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/dashboard.html")
	if err != nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

// GenerateQRCode creates a scannable PNG image linking to the customer check-in page
func (h *ShopkeeperHandler) GenerateQRCode(w http.ResponseWriter, r *http.Request) {
	shopID := r.PathValue("shopID")
	if shopID == "" {
		http.Error(w, "missing shop ID", http.StatusBadRequest)
		return
	}

	// Construct the URL the customer will scan
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	customerURL := fmt.Sprintf("%s%s/c/%s", scheme, r.Host, shopID)

	// Generate the QR code as a PNG
	png, err := qrcode.Encode(customerURL, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "failed to generate QR", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

// RedeemPrize processes a claimed reward
func (h *ShopkeeperHandler) RedeemPrize(w http.ResponseWriter, r *http.Request) {
	shopID := r.PathValue("shopID")
	customerID := r.PathValue("customerID")

	if err := h.shopCustomerRepo.RedeemReward(r.Context(), shopID, customerID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}
