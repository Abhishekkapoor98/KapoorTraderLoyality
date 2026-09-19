package main

import (
	"context"
	"log"
	"net/http"
	"os"

	loyalty "kapoortrader-loyalty"
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

	// Execute embedded migrations automatically
	initSQL, err := loyalty.MigrationsFS.ReadFile("migrations/001_init.sql")
	if err == nil {
		if _, err := db.ExecContext(ctx, string(initSQL)); err != nil {
			log.Printf("Init migration skipped or already applied: %v\n", err)
		}
	}

	seedSQL, err := loyalty.MigrationsFS.ReadFile("migrations/002_seed.sql")
	if err == nil {
		if _, err := db.ExecContext(ctx, string(seedSQL)); err != nil {
			log.Printf("Seed migration skipped or already applied: %v\n", err)
		}
	}

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

	// Go 1.22+ wildcard routing format
	http.HandleFunc("GET /c/{shopID}", customerHandler.ServeCustomerLoginPage)
	http.HandleFunc("POST /c/{shopID}/login", customerHandler.LoginAndPurchase)

	// Shopkeeper Routes
	http.HandleFunc("GET /shopkeeper/{shopID}/dashboard", shopkeeperHandler.ServeDashboardPage)
	http.HandleFunc("GET /shopkeeper/{shopID}/qr", shopkeeperHandler.GenerateQRCode)
	http.HandleFunc("POST /shopkeeper/{shopID}/redeem/{customerID}", shopkeeperHandler.RedeemPrize)

	// 1. Add a root route to prevent the 404 Not Found error
	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// Prevent wildcard catching everything if you only want it on the exact root
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte("Kapoor Trader Loyalty App is live!"))
	})

	// 2. Read Render's dynamic PORT environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback for local development
	}

	log.Println("server running on port " + port)

	// 3. Listen on the dynamic port
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
