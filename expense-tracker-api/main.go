package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"expense-tracker-api/config"
	categoryrepo "expense-tracker-api/pkg/categories/repository"
	categorhandler "expense-tracker-api/pkg/categories/handler"
	familyrepo "expense-tracker-api/pkg/families/repository"
	familyhandler "expense-tracker-api/pkg/families/handler"
	"expense-tracker-api/pkg/handler"
	paymentmethodrepo "expense-tracker-api/pkg/paymentmethods/repository"
	paymentmethodhandler "expense-tracker-api/pkg/paymentmethods/handler"
	"expense-tracker-api/pkg/response"
	"expense-tracker-api/pkg/transactions/dashboard"
	transactionrepo "expense-tracker-api/pkg/transactions/repository"
	transactionhandler "expense-tracker-api/pkg/transactions/handler"
	upiapprepo "expense-tracker-api/pkg/upiapps/repository"
	upiapphandler "expense-tracker-api/pkg/upiapps/handler"
	userrepo "expense-tracker-api/pkg/users/repository"
	"expense-tracker-api/pkg/log"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"expense-tracker-api/internal/middleware"
)

func main() {
	cfg := config.LoadConfig()
	log.Init()

	pool := cfg.DB

	userRepo := userrepo.NewUserRepository(pool)
	familyRepo := familyrepo.NewFamilyRepository(pool)
	transactionRepo := transactionrepo.NewTransactionRepository(pool)
	categoryRepo := categoryrepo.NewCategoryRepository(pool)
	paymentMethodRepo := paymentmethodrepo.NewPaymentMethodRepository(pool)
	upiAppRepo := upiapprepo.NewUPIAppRepository(pool)

	userHandler := handler.NewUserHandler(pool)
	familyHandler := familyhandler.NewFamilyHandler(familyRepo)
	transactionHandler := transactionhandler.NewTransactionHandler(transactionRepo)
	categoryHandler := categorhandler.NewCategoryHandler(categoryRepo)
	paymentMethodHandler := paymentmethodhandler.NewPaymentMethodHandler(paymentMethodRepo)
	upiAppHandler := upiapphandler.NewUPIAppHandler(upiAppRepo)
	dashboardHandler := dashboard.NewDashboardHandler(pool)

	authMiddleware := middleware.NewAuthMiddleware(pool)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	r.Use(corsMiddleware)
	r.Use(authMiddleware.Middleware)

	// Health check (no auth required)
	r.Get("/health", healthCheck(pool))

	r.Get("/", homeHandler)

	// Auth routes
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/google", userHandler.RegisterGoogle)
		r.Post("/verify-otp", userHandler.VerifyOTP)
		r.Post("/refresh", userHandler.RefreshToken)
		r.Post("/logout", userHandler.Logout)
	})

	// User routes
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Get("/me", makeHandler(userHandler.GetMe))
		r.Patch("/me", makeHandler(userHandler.UpdateMe))
	})

	// Family routes
	r.Route("/api/v1/families", func(r chi.Router) {
		r.Post("/", makeHandler(familyHandler.CreateFamily))
		r.Get("/", makeHandler(familyHandler.GetFamily))
		r.Get("/me", makeHandler(familyHandler.ListMembers))
		r.Post("/invite", makeHandler(familyHandler.InviteMember))
	})

	// Dashboard - family head only
	r.Route("/api/v1/dashboard", func(r chi.Router) {
		r.Get("/", makeHandler(dashboardHandler.GetDashboard))
	})

	// Transaction routes
	r.Route("/api/v1/transactions", func(r chi.Router) {
		r.Post("/", makeHandler(transactionHandler.CreateTransaction))
		r.Get("/", makeHandler(transactionHandler.ListTransactions))
		r.Get("/me", makeHandler(transactionHandler.ListMyTransactions))
	})

	// Category routes
	r.Route("/api/v1/categories", func(r chi.Router) {
		r.Get("/", makeHandler(categoryHandler.ListByType))
	})

	// Payment method routes
	r.Route("/api/v1/payment-methods", func(r chi.Router) {
		r.Get("/", makeHandler(paymentMethodHandler.ListAll))
	})

	// UPI app routes
	r.Route("/api/v1/upi-apps", func(r chi.Router) {
		r.Get("/", makeHandler(upiAppHandler.ListAll))
	})

	_ = userRepo // Used indirectly through handlers

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Std.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Std.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Std.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Std.Fatalf("Server forced to shutdown: %v", err)
	}

	pool.Close()
	log.Std.Println("Server exited gracefully")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Data: map[string]string{
			"message": "Welcome to Expense Tracker API v1",
		},
	})
}

func makeHandler(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}
}

func healthCheck(pool interface{ Ping(context.Context) error }) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		dbStatus := "ok"
		redisStatus := "ok"

		if err := pool.Ping(ctx); err != nil {
			dbStatus = "error: " + err.Error()
		}

		if cfg := config.LoadConfig(); cfg.Rdb != nil {
			if err := cfg.Rdb.Ping(ctx).Err(); err != nil {
				redisStatus = "error: " + err.Error()
			}
		}

		status := http.StatusOK
		if dbStatus != "ok" || redisStatus != "ok" {
			status = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "ok",
			"db":     dbStatus,
			"redis":  redisStatus,
		})
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
