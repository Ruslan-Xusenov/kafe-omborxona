package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"kafe-omborxona/internal/config"
	"kafe-omborxona/internal/handler"
	mw "kafe-omborxona/internal/middleware"
	"kafe-omborxona/internal/repository"
	"kafe-omborxona/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// Database
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("DB ping error: %v", err)
	}
	log.Println("✅ Database connected")

	// Run migrations
	runMigrations(pool)

	// Seed admin
	seedAdmin(pool)

	// Repositories
	userRepo := repository.NewUserRepo(pool)
	categoryRepo := repository.NewCategoryRepo(pool)
	supplierRepo := repository.NewSupplierRepo(pool)
	productRepo := repository.NewProductRepo(pool)
	txnRepo := repository.NewTransactionRepo(pool)
	recipeRepo := repository.NewRecipeRepo(pool)
	debtRepo := repository.NewDebtRepo(pool)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	tgSvc, _ := service.NewTelegramService(cfg.TelegramBotToken)

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	catH := handler.NewCategoryHandler(categoryRepo)
	supH := handler.NewSupplierHandler(supplierRepo)
	prodH := handler.NewProductHandler(productRepo)
	txnH := handler.NewTransactionHandler(txnRepo, tgSvc)
	dashH := handler.NewDashboardHandler(txnRepo, debtRepo, tgSvc)
	recH := handler.NewRecipeHandler(recipeRepo)
	debtH := handler.NewDebtHandler(debtRepo)
	exportH := handler.NewExportHandler(txnRepo)

	// Router
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("POST /api/auth/login", authH.Login)

	// Protected routes - wrap with auth middleware
	authMw := mw.Auth(authSvc)
	adminMw := mw.RequireRole("admin")

	// Auth
	mux.Handle("GET /api/auth/me", authMw(http.HandlerFunc(authH.Me)))

	// Users (admin only)
	mux.Handle("GET /api/users", authMw(adminMw(http.HandlerFunc(authH.GetUsers))))
	mux.Handle("POST /api/users", authMw(adminMw(http.HandlerFunc(authH.CreateUser))))
	mux.Handle("PUT /api/users/{id}", authMw(adminMw(http.HandlerFunc(authH.UpdateUser))))
	mux.Handle("DELETE /api/users/{id}", authMw(adminMw(http.HandlerFunc(authH.DeleteUser))))

	// Categories
	mux.Handle("GET /api/categories", authMw(http.HandlerFunc(catH.GetAll)))
	mux.Handle("POST /api/categories", authMw(http.HandlerFunc(catH.Create)))
	mux.Handle("PUT /api/categories/{id}", authMw(http.HandlerFunc(catH.Update)))
	mux.Handle("DELETE /api/categories/{id}", authMw(http.HandlerFunc(catH.Delete)))

	// Suppliers
	mux.Handle("GET /api/suppliers", authMw(http.HandlerFunc(supH.GetAll)))
	mux.Handle("POST /api/suppliers", authMw(http.HandlerFunc(supH.Create)))
	mux.Handle("PUT /api/suppliers/{id}", authMw(http.HandlerFunc(supH.Update)))
	mux.Handle("DELETE /api/suppliers/{id}", authMw(http.HandlerFunc(supH.Delete)))

	// Products
	mux.Handle("GET /api/products", authMw(http.HandlerFunc(prodH.GetAll)))
	mux.Handle("POST /api/products", authMw(http.HandlerFunc(prodH.Create)))
	mux.Handle("PUT /api/products/{id}", authMw(http.HandlerFunc(prodH.Update)))
	mux.Handle("DELETE /api/products/{id}", authMw(http.HandlerFunc(prodH.Delete)))

	// Transactions
	mux.Handle("GET /api/transactions", authMw(http.HandlerFunc(txnH.GetAll)))
	mux.Handle("POST /api/transactions", authMw(http.HandlerFunc(txnH.Create)))
	mux.Handle("DELETE /api/transactions/{id}", authMw(adminMw(http.HandlerFunc(txnH.Delete))))
	mux.Handle("GET /api/export/transactions", authMw(http.HandlerFunc(exportH.ExportTransactions)))

	// Recipes
	mux.Handle("GET /api/recipes", authMw(http.HandlerFunc(recH.GetAll)))
	mux.Handle("POST /api/recipes", authMw(http.HandlerFunc(recH.Create)))
	mux.Handle("PUT /api/recipes/{id}", authMw(http.HandlerFunc(recH.Update)))
	mux.Handle("DELETE /api/recipes/{id}", authMw(http.HandlerFunc(recH.Delete)))

	// Debts
	mux.Handle("GET /api/debts", authMw(http.HandlerFunc(debtH.GetAll)))
	mux.Handle("POST /api/debts", authMw(http.HandlerFunc(debtH.Create)))
	mux.Handle("POST /api/debts/{id}/pay", authMw(http.HandlerFunc(debtH.Pay)))
	mux.Handle("DELETE /api/debts/{id}", authMw(adminMw(http.HandlerFunc(debtH.Delete))))

	// Dashboard
	mux.Handle("GET /api/dashboard/summary", authMw(http.HandlerFunc(dashH.Summary)))
	mux.Handle("GET /api/dashboard/inventory", authMw(http.HandlerFunc(dashH.Inventory)))
	mux.Handle("GET /api/dashboard/alerts", authMw(http.HandlerFunc(dashH.Alerts)))
	mux.Handle("GET /api/dashboard/profit", authMw(adminMw(http.HandlerFunc(dashH.Profit))))
	mux.Handle("GET /api/dashboard/top-products", authMw(adminMw(http.HandlerFunc(dashH.TopProducts))))
	mux.Handle("POST /api/dashboard/trigger-report", authMw(adminMw(http.HandlerFunc(dashH.TriggerReport))))

	// CORS wrapper
	corsHandler := mw.CORS(cfg.FrontendURL)(mux)

	// Server
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      corsHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🚀 Server started on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func runMigrations(pool *pgxpool.Pool) {
	migrations := []string{
		"001_init.up.sql", 
		"002_add_minstock_barcode_expiry.up.sql", 
		"003_add_recipes.up.sql",
		"004_add_debts.up.sql",
		"005_cascade_deletes.up.sql",
		"006_cascade_recipe_ingredients.up.sql",
	}
	for _, m := range migrations {
		migration, err := os.ReadFile("migrations/" + m)
		if err != nil {
			log.Printf("⚠️  Migration file %s not found, skipping: %v", m, err)
			continue
		}
		_, err = pool.Exec(context.Background(), string(migration))
		if err != nil {
			log.Printf("⚠️  Migration %s may already be applied: %v", m, err)
		} else {
			log.Printf("✅ Migration applied: %s", m)
		}
	}
}

func seedAdmin(pool *pgxpool.Pool) {
	var count int
	pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM users WHERE role='admin'`).Scan(&count)
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("⚠️  Failed to hash admin password: %v", err)
		return
	}

	_, err = pool.Exec(context.Background(),
		`INSERT INTO users (username, password_hash, full_name, role) VALUES ($1, $2, $3, $4)`,
		"admin", string(hash), "Administrator", "admin")
	if err != nil {
		log.Printf("⚠️  Failed to seed admin: %v", err)
	} else {
		fmt.Println("✅ Admin user created (admin / admin123)")
	}
}
