package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"business-mgmt-backend/internal/auditlog"
	"business-mgmt-backend/internal/auth"
	"business-mgmt-backend/internal/configuration"
	"business-mgmt-backend/internal/customer"
	"business-mgmt-backend/internal/expense"
	"business-mgmt-backend/internal/file"
	"business-mgmt-backend/internal/inventory"
	"business-mgmt-backend/internal/notification"
	"business-mgmt-backend/internal/notificationtemplate"
	"business-mgmt-backend/internal/payment"
	"business-mgmt-backend/internal/product"
	"business-mgmt-backend/internal/productcategory"
	"business-mgmt-backend/internal/purchase"
	"business-mgmt-backend/internal/rbac"
	"business-mgmt-backend/internal/reminderrule"
	"business-mgmt-backend/internal/report"
	"business-mgmt-backend/internal/sale"
	"business-mgmt-backend/internal/store"
	"business-mgmt-backend/internal/supplier"
	"business-mgmt-backend/internal/system"
	"business-mgmt-backend/internal/user"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// 1. System & Health
	systemHandler := system.NewHandler()
	system.RegisterRoutes(mux, systemHandler, auth.AuthMiddleware)

	// 2. Auth
	authUserRepo := NewStubAuthUserRepo()
	authSessionRepo := NewStubSessionRepo()
	authService := auth.NewService(authUserRepo, authSessionRepo)
	authHandler := auth.NewHandler(authService)
	auth.RegisterRoutes(mux, authHandler)

	// 3. User Management
	userHandler := user.NewHandler(user.NewService(&StubUserRepo{}))
	user.RegisterRoutes(mux, userHandler, auth.AuthMiddleware)

	// 4. RBAC
	rbacHandler := rbac.NewHandler(rbac.NewService(&StubRBACRepo{}))
	rbac.RegisterRoutes(mux, rbacHandler, auth.AuthMiddleware)

	// 5. Stores
	storeHandler := store.NewHandler(store.NewService(&StubStoreRepo{}))
	store.RegisterRoutes(mux, storeHandler, auth.AuthMiddleware)

	// 6. Customers
	customerHandler := customer.NewHandler(customer.NewService(&StubCustomerRepo{}))
	customer.RegisterRoutes(mux, customerHandler, auth.AuthMiddleware)

	// 7. Product Categories
	productCategoryHandler := productcategory.NewHandler(productcategory.NewService(&StubProductCategoryRepo{}))
	productcategory.RegisterRoutes(mux, productCategoryHandler, auth.AuthMiddleware)

	// 8. Products
	productHandler := product.NewHandler(product.NewService(&StubProductRepo{}))
	product.RegisterRoutes(mux, productHandler, auth.AuthMiddleware)

	// 9. Suppliers
	supplierHandler := supplier.NewHandler(supplier.NewService(&StubSupplierRepo{}))
	supplier.RegisterRoutes(mux, supplierHandler, auth.AuthMiddleware)

	// 10. Purchases
	purchaseHandler := purchase.NewHandler(purchase.NewService(&StubPurchaseRepo{}))
	purchase.RegisterRoutes(mux, purchaseHandler, auth.AuthMiddleware)

	// 11. Inventory
	inventoryHandler := inventory.NewHandler(inventory.NewService(&StubInventoryRepo{}))
	inventory.RegisterRoutes(mux, inventoryHandler, auth.AuthMiddleware)

	// 12. Sales
	saleHandler := sale.NewHandler(sale.NewService(&StubSaleRepo{}, &StubProductInventoryFetcher{}))
	sale.RegisterRoutes(mux, saleHandler, auth.AuthMiddleware)

	// 13. Payments
	paymentHandler := payment.NewHandler(payment.NewService(&StubPaymentRepo{}))
	payment.RegisterRoutes(mux, paymentHandler, auth.AuthMiddleware)

	// 14. Expenses
	expenseHandler := expense.NewHandler(expense.NewService(&StubExpenseRepo{}))
	expense.RegisterRoutes(mux, expenseHandler, auth.AuthMiddleware)

	// 15. Configuration
	configHandler := configuration.NewHandler(configuration.NewService(&StubConfigRepo{}))
	configuration.RegisterRoutes(mux, configHandler, auth.AuthMiddleware)

	// 16. Reminder Rules
	reminderRuleHandler := reminderrule.NewHandler(reminderrule.NewService(&StubReminderRuleRepo{}))
	reminderrule.RegisterRoutes(mux, reminderRuleHandler, auth.AuthMiddleware)

	// 17. Notification Templates
	notificationTemplateHandler := notificationtemplate.NewHandler(notificationtemplate.NewService(&StubNotificationTemplateRepo{}))
	notificationtemplate.RegisterRoutes(mux, notificationTemplateHandler, auth.AuthMiddleware)

	// 18. Notifications
	notificationHandler := notification.NewHandler(notification.NewService(&StubNotificationRepo{}, &StubNotificationValidator{}))
	notification.RegisterRoutes(mux, notificationHandler, auth.AuthMiddleware)

	// 19. Reports
	reportHandler := report.NewHandler(report.NewService(&StubReportRepo{}))
	report.RegisterRoutes(mux, reportHandler, auth.AuthMiddleware)

	// 20. Audit Logs
	auditLogHandler := auditlog.NewHandler(auditlog.NewService(&StubAuditLogRepo{}))
	auditlog.RegisterRoutes(mux, auditLogHandler, auth.AuthMiddleware)

	// 21. Files
	fileHandler := file.NewHandler(file.NewService())
	file.RegisterRoutes(mux, fileHandler, auth.AuthMiddleware)

	// Wrap root with logging and CORS
	rootHandler := loggingMiddleware(mux)

	fmt.Println("==================================================================")
	fmt.Println("🚀 Business Management Software API Server")
	fmt.Printf("🌐 Server listening on http://localhost:%s\n", port)
	fmt.Println("🔑 Default Admin Credentials for Postman Testing:")
	fmt.Printf("   Email:    %s\n", SeedEmail)
	fmt.Printf("   Password: %s\n", SeedPass)
	fmt.Println("📦 All 21 modules registered successfully!")
	fmt.Println("==================================================================")

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      rootHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Idempotency-Key")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s (%v)", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}
