package main

import (
	"context"
	"log"
	"net/http"
	"os"

	deliveryHttp "kapoortrader-loyalty/internal/delivery/http"
	"kapoortrader-loyalty/internal/infrastructure/postgres"
	"kapoortrader-loyalty/internal/usecase"
)

func main() {
	ctx := context.Background()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://loyalty:loyalty@localhost:5432/loyalty?sslmode=disable"
	}

	db, err := postgres.NewDB(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. Initialize Repositories
	customerRepo := postgres.NewCustomerRepository(db)
	shopCustomerRepo := postgres.NewShopCustomerRepository(db)

	// 2. Initialize Use Cases (Service Layer)
	customerService := usecase.NewCustomerService(customerRepo, shopCustomerRepo)

	// 3. Initialize HTTP Handlers
	customerHandler := deliveryHttp.NewCustomerHandler(customerService)

	shopkeeperHandler := deliveryHttp.NewShopkeeperHandler(shopCustomerRepo)

	http.HandleFunc("GET /shopkeeper/{shopID}/dashboard/data", shopkeeperHandler.GetDashboardData)

	// 4. Register Routes
	http.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			http.Error(w, "database unavailable", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Go 1.22+ wildcard routing format[cite: 1]
	http.HandleFunc("GET /c/{shopID}", customerHandler.ServeCustomerLoginPage)
	http.HandleFunc("POST /c/{shopID}/login", customerHandler.LoginAndPurchase)

	// Shopkeeper Routes
	http.HandleFunc("GET /shopkeeper/{shopID}/dashboard", shopkeeperHandler.ServeDashboardPage)

	// Add these two new routes:
	http.HandleFunc("GET /shopkeeper/{shopID}/qr", shopkeeperHandler.GenerateQRCode)
	http.HandleFunc("POST /shopkeeper/{shopID}/redeem/{customerID}", shopkeeperHandler.RedeemPrize)

	log.Println("server running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
