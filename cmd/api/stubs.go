package main

import (
	"context"
	"errors"
	"time"

	"business-mgmt-backend/internal/models"
	"business-mgmt-backend/internal/payment"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrNotFound = errors.New("resource not found")

// Seed admin user credentials
var (
	SeedUserID = uuid.New()
	SeedEmail  = "admin@example.com"
	SeedPass   = "Password123!"
)

// --- Auth Repositories ---

type StubAuthUserRepo struct {
	adminUser *models.User
}

func NewStubAuthUserRepo() *StubAuthUserRepo {
	hash, _ := bcrypt.GenerateFromPassword([]byte(SeedPass), bcrypt.DefaultCost)
	now := time.Now()
	return &StubAuthUserRepo{
		adminUser: &models.User{
			ID:           SeedUserID,
			Email:        SeedEmail,
			PasswordHash: string(hash),
			FullName:     "System Administrator",
			Status:       "ACTIVE",
			Roles:        []string{"ADMIN"},
			Permissions:  []string{"*"},
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}
}

func (r *StubAuthUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == r.adminUser.Email {
		return r.adminUser, nil
	}
	return nil, ErrNotFound
}

func (r *StubAuthUserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	if id == r.adminUser.ID {
		return r.adminUser, nil
	}
	return nil, ErrNotFound
}

func (r *StubAuthUserRepo) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	r.adminUser.LastLoginAt = &now
	return nil
}

func (r *StubAuthUserRepo) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	if id == r.adminUser.ID {
		r.adminUser.PasswordHash = passwordHash
		return nil
	}
	return ErrNotFound
}

type StubSessionRepo struct {
	sessions map[string]*models.UserSession
}

func NewStubSessionRepo() *StubSessionRepo {
	return &StubSessionRepo{sessions: make(map[string]*models.UserSession)}
}

func (r *StubSessionRepo) CreateSession(ctx context.Context, s *models.UserSession) error {
	r.sessions[s.RefreshTokenHash] = s
	return nil
}

func (r *StubSessionRepo) GetSessionByTokenHash(ctx context.Context, hash string) (*models.UserSession, error) {
	if s, ok := r.sessions[hash]; ok {
		return s, nil
	}
	return nil, ErrNotFound
}

func (r *StubSessionRepo) RevokeSession(ctx context.Context, id uuid.UUID) error {
	return nil
}

// --- User Repo ---

type StubUserRepo struct{}

func (r *StubUserRepo) CreateUser(ctx context.Context, user *models.User) error { return nil }
func (r *StubUserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return &models.User{ID: id, FullName: "Test User", Email: "test@example.com", Status: "ACTIVE"}, nil
}
func (r *StubUserRepo) ListUsers(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.User, int, error) {
	return []models.User{}, 0, nil
}
func (r *StubUserRepo) UpdateUser(ctx context.Context, id uuid.UUID, fullName, phone *string) error { return nil }
func (r *StubUserRepo) UpdateUserStatus(ctx context.Context, id uuid.UUID, status string) error   { return nil }
func (r *StubUserRepo) AssignRoles(ctx context.Context, id uuid.UUID, roleCodes []string) error    { return nil }
func (r *StubUserRepo) AssignStores(ctx context.Context, id uuid.UUID, storeIDs []uuid.UUID) error { return nil }

// --- RBAC Repo ---

type StubRBACRepo struct{}

func (r *StubRBACRepo) ListRoles(ctx context.Context) ([]models.Role, error) {
	return []models.Role{
		{ID: uuid.New(), Code: "ADMIN", Name: "Administrator"},
		{ID: uuid.New(), Code: "STORE_MANAGER", Name: "Store Manager"},
	}, nil
}
func (r *StubRBACRepo) GetRoleByID(ctx context.Context, roleID uuid.UUID) (*models.Role, error) {
	return &models.Role{ID: roleID, Code: "ADMIN", Name: "Administrator"}, nil
}
func (r *StubRBACRepo) ListPermissions(ctx context.Context) ([]models.Permission, error) {
	return []models.Permission{
		{ID: uuid.New(), Code: "users:create"},
		{ID: uuid.New(), Code: "sales:create"},
	}, nil
}
func (r *StubRBACRepo) UpdateRolePermissions(ctx context.Context, roleID uuid.UUID, permissionCodes []string) error {
	return nil
}

// --- Store Repo ---

type StubStoreRepo struct{}

func (r *StubStoreRepo) CreateStore(ctx context.Context, store *models.StoreModel) error { return nil }
func (r *StubStoreRepo) GetStoreByID(ctx context.Context, id uuid.UUID) (*models.StoreModel, error) {
	return &models.StoreModel{ID: id, Name: "Main Store", Code: "STORE-01", Status: "ACTIVE"}, nil
}
func (r *StubStoreRepo) ListStores(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.StoreModel, int, error) {
	return []models.StoreModel{{ID: uuid.New(), Name: "Main Branch", Code: "MAIN-01", Status: "ACTIVE"}}, 1, nil
}
func (r *StubStoreRepo) UpdateStore(ctx context.Context, id uuid.UUID, store *models.StoreModel) error { return nil }
func (r *StubStoreRepo) UpdateStoreStatus(ctx context.Context, id uuid.UUID, status string) error     { return nil }

// --- Customer Repo ---

type StubCustomerRepo struct{}

func (r *StubCustomerRepo) CreateCustomer(ctx context.Context, customer *models.CustomerModel) error { return nil }
func (r *StubCustomerRepo) GetCustomerByID(ctx context.Context, id uuid.UUID) (*models.CustomerModel, error) {
	return &models.CustomerModel{ID: id, Name: "Jane Doe", PhoneNumber: "9876543210", Status: "ACTIVE"}, nil
}
func (r *StubCustomerRepo) GetCustomerByPhone(ctx context.Context, phone string) (*models.CustomerModel, error) {
	return nil, ErrNotFound
}
func (r *StubCustomerRepo) ListCustomers(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.CustomerModel, int, error) {
	return []models.CustomerModel{}, 0, nil
}
func (r *StubCustomerRepo) UpdateCustomer(ctx context.Context, id uuid.UUID, customer *models.CustomerModel) error { return nil }
func (r *StubCustomerRepo) UpdateCustomerStatus(ctx context.Context, id uuid.UUID, status string) error           { return nil }
func (r *StubCustomerRepo) LinkCustomerToStore(ctx context.Context, customerID, storeID uuid.UUID, addedBy *uuid.UUID) error { return nil }
func (r *StubCustomerRepo) GetCustomerStores(ctx context.Context, customerID uuid.UUID) ([]models.StoreModel, error) { return []models.StoreModel{}, nil }

// --- Product Category Repo ---

type StubProductCategoryRepo struct{}

func (r *StubProductCategoryRepo) CreateCategory(ctx context.Context, category *models.ProductCategory) error { return nil }
func (r *StubProductCategoryRepo) GetCategoryByID(ctx context.Context, id uuid.UUID) (*models.ProductCategory, error) {
	return &models.ProductCategory{ID: id, Name: "General", Status: "ACTIVE"}, nil
}
func (r *StubProductCategoryRepo) GetCategoryByName(ctx context.Context, name string) (*models.ProductCategory, error) {
	return nil, ErrNotFound
}
func (r *StubProductCategoryRepo) ListCategories(ctx context.Context) ([]models.ProductCategory, error) {
	return []models.ProductCategory{{ID: uuid.New(), Name: "Electronics", Status: "ACTIVE"}}, nil
}
func (r *StubProductCategoryRepo) UpdateCategory(ctx context.Context, id uuid.UUID, category *models.ProductCategory) error { return nil }
func (r *StubProductCategoryRepo) UpdateCategoryStatus(ctx context.Context, id uuid.UUID, status string) error             { return nil }

// --- Product Repo ---

type StubProductRepo struct{}

func (r *StubProductRepo) CreateProduct(ctx context.Context, product *models.ProductModel) error { return nil }
func (r *StubProductRepo) GetProductByID(ctx context.Context, id uuid.UUID) (*models.ProductModel, error) {
	return &models.ProductModel{ID: id, Name: "Sample Item", SKU: "SMP-001", Status: "ACTIVE", DefaultSellingPrice: "100.00"}, nil
}
func (r *StubProductRepo) GetProductBySKU(ctx context.Context, sku string) (*models.ProductModel, error) {
	return nil, ErrNotFound
}
func (r *StubProductRepo) ListProducts(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.ProductModel, int, error) {
	return []models.ProductModel{}, 0, nil
}
func (r *StubProductRepo) UpdateProduct(ctx context.Context, id uuid.UUID, product *models.ProductModel) error { return nil }
func (r *StubProductRepo) UpdateProductStatus(ctx context.Context, id uuid.UUID, status string) error         { return nil }
func (r *StubProductRepo) GetProductInventoryAcrossStores(ctx context.Context, id uuid.UUID) ([]models.ProductInventoryData, error) {
	return []models.ProductInventoryData{}, nil
}

// --- Supplier Repo ---

type StubSupplierRepo struct{}

func (r *StubSupplierRepo) CreateSupplier(ctx context.Context, supplier *models.SupplierModel) error { return nil }
func (r *StubSupplierRepo) GetSupplierByID(ctx context.Context, id uuid.UUID) (*models.SupplierModel, error) {
	return &models.SupplierModel{ID: id, Name: "Acme Supplies", Status: "ACTIVE"}, nil
}
func (r *StubSupplierRepo) ListSuppliers(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.SupplierModel, int, error) {
	return []models.SupplierModel{}, 0, nil
}
func (r *StubSupplierRepo) UpdateSupplier(ctx context.Context, id uuid.UUID, supplier *models.SupplierModel) error { return nil }
func (r *StubSupplierRepo) UpdateSupplierStatus(ctx context.Context, id uuid.UUID, status string) error           { return nil }

// --- Purchase Repo ---

type StubPurchaseRepo struct{}

func (r *StubPurchaseRepo) CreatePurchase(ctx context.Context, purchase *models.PurchaseModel) error { return nil }
func (r *StubPurchaseRepo) GetPurchaseByID(ctx context.Context, id uuid.UUID) (*models.PurchaseModel, error) {
	return &models.PurchaseModel{ID: id, PurchaseNumber: "PO-001", Status: "DRAFT"}, nil
}
func (r *StubPurchaseRepo) ListPurchases(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.PurchaseModel, int, error) {
	return []models.PurchaseModel{}, 0, nil
}
func (r *StubPurchaseRepo) UpdatePurchase(ctx context.Context, purchase *models.PurchaseModel) error { return nil }
func (r *StubPurchaseRepo) UpdatePurchaseStatus(ctx context.Context, id uuid.UUID, status string) error { return nil }
func (r *StubPurchaseRepo) ExecuteApprovalTransaction(ctx context.Context, purchaseID uuid.UUID, approvedBy *uuid.UUID, fn func(context.Context, *models.PurchaseModel) error) error {
	return nil
}

// --- Inventory Repo ---

type StubInventoryRepo struct{}

func (r *StubInventoryRepo) GetStoreInventory(ctx context.Context, storeID uuid.UUID, page, pageSize int, filters map[string]string) ([]models.InventoryModel, int, error) {
	return []models.InventoryModel{}, 0, nil
}
func (r *StubInventoryRepo) GetProductInventoryInStore(ctx context.Context, storeID, productID uuid.UUID) (*models.InventoryModel, error) {
	return &models.InventoryModel{StoreID: storeID, ProductID: productID, AvailableQuantity: "50"}, nil
}
func (r *StubInventoryRepo) GetLowStockReport(ctx context.Context, storeID uuid.UUID, threshold string, page, pageSize int) ([]models.InventoryModel, int, error) {
	return []models.InventoryModel{}, 0, nil
}
func (r *StubInventoryRepo) GetInventoryMovements(ctx context.Context, storeID, productID uuid.UUID, page, pageSize int, filters map[string]string) ([]models.InventoryMovementModel, int, error) {
	return []models.InventoryMovementModel{}, 0, nil
}
func (r *StubInventoryRepo) ExecuteInventoryAdjustment(ctx context.Context, storeID uuid.UUID, adjustment models.CreateInventoryAdjustmentRequest, createdBy uuid.UUID, fn func(context.Context, *models.InventoryModel) error) error {
	return nil
}

// --- Sale Repo & Fetcher ---

type StubSaleRepo struct{}

func (r *StubSaleRepo) GetInvoiceByID(ctx context.Context, id uuid.UUID) (*models.InvoiceModel, error) {
	return &models.InvoiceModel{ID: id, InvoiceNumber: "INV-001", PaymentStatus: "PAID", TotalAmount: "100.00"}, nil
}
func (r *StubSaleRepo) GetInvoiceByNumber(ctx context.Context, number string) (*models.InvoiceModel, error) {
	return &models.InvoiceModel{ID: uuid.New(), InvoiceNumber: number, PaymentStatus: "PAID"}, nil
}
func (r *StubSaleRepo) ListInvoices(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.InvoiceModel, int, error) {
	return []models.InvoiceModel{}, 0, nil
}
func (r *StubSaleRepo) GetInvoicePayments(ctx context.Context, invoiceID uuid.UUID) ([]models.PaymentAllocationModel, error) {
	return []models.PaymentAllocationModel{}, nil
}
func (r *StubSaleRepo) ExecuteSaleTransaction(ctx context.Context, invoice *models.InvoiceModel, initialPayment *models.InitialPaymentModel) error {
	return nil
}
func (r *StubSaleRepo) ExecuteCancelSaleTransaction(ctx context.Context, invoiceID uuid.UUID, reason string) error {
	return nil
}

type StubProductInventoryFetcher struct{}

func (f *StubProductInventoryFetcher) GetProductGSTRate(ctx context.Context, productID uuid.UUID) (float64, error) {
	return 18.0, nil
}
func (f *StubProductInventoryFetcher) GetAverageUnitCost(ctx context.Context, storeID, productID uuid.UUID) (float64, error) {
	return 50.0, nil
}

// --- Payment Repo ---

type StubPaymentRepo struct{}

func (r *StubPaymentRepo) GetPaymentByID(ctx context.Context, id uuid.UUID) (*models.PaymentModel, error) {
	return &models.PaymentModel{ID: id, Amount: "100.00", PaymentMethod: "CASH", Status: "COMPLETED"}, nil
}
func (r *StubPaymentRepo) ListPayments(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.PaymentModel, int, error) {
	return []models.PaymentModel{}, 0, nil
}
func (r *StubPaymentRepo) ExecutePaymentTransaction(ctx context.Context, p *models.PaymentModel, strategy string, allocs []models.CreatePaymentAllocationStruct, fn func(context.Context, payment.InvoiceLocker) ([]models.PaymentAllocationModel, error)) error {
	return nil
}
func (r *StubPaymentRepo) ExecutePaymentReversalTransaction(ctx context.Context, paymentID uuid.UUID, reason string, reversedBy uuid.UUID) error {
	return nil
}

// --- Expense Repo ---

type StubExpenseRepo struct{}

func (r *StubExpenseRepo) CreateExpenseCategory(ctx context.Context, category *models.ExpenseCategoryModel) error { return nil }
func (r *StubExpenseRepo) GetExpenseCategoryByID(ctx context.Context, id uuid.UUID) (*models.ExpenseCategoryModel, error) {
	return &models.ExpenseCategoryModel{ID: id, Name: "Rent"}, nil
}
func (r *StubExpenseRepo) ListExpenseCategories(ctx context.Context, filters map[string]string) ([]models.ExpenseCategoryModel, int, error) {
	return []models.ExpenseCategoryModel{}, 0, nil
}
func (r *StubExpenseRepo) UpdateExpenseCategory(ctx context.Context, category *models.ExpenseCategoryModel) error { return nil }
func (r *StubExpenseRepo) CreateExpense(ctx context.Context, expense *models.ExpenseModel) error                 { return nil }
func (r *StubExpenseRepo) GetExpenseByID(ctx context.Context, id uuid.UUID) (*models.ExpenseModel, error) {
	return &models.ExpenseModel{ID: id, Amount: "500.00", Description: "Office Supplies"}, nil
}
func (r *StubExpenseRepo) ListExpenses(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.ExpenseModel, int, error) {
	return []models.ExpenseModel{}, 0, nil
}
func (r *StubExpenseRepo) UpdateExpense(ctx context.Context, expense *models.ExpenseModel) error { return nil }
func (r *StubExpenseRepo) ExecuteExpenseReversalTransaction(ctx context.Context, expenseID uuid.UUID, reason string, currentUserID uuid.UUID) error {
	return nil
}

// --- Configuration Repo ---

type StubConfigRepo struct{}

func (r *StubConfigRepo) GetCompanySettings(ctx context.Context) (*models.CompanySettingsModel, error) {
	return &models.CompanySettingsModel{
		CompanyName: "Test Business Corp",
		Currency:    "INR",
		Timezone:    "Asia/Kolkata",
	}, nil
}
func (r *StubConfigRepo) UpdateCompanySettings(ctx context.Context, settings *models.CompanySettingsModel) error {
	return nil
}

// --- Reminder Rule Repo ---

type StubReminderRuleRepo struct{}

func (r *StubReminderRuleRepo) CreateRule(ctx context.Context, rule *models.ReminderRuleModel) error { return nil }
func (r *StubReminderRuleRepo) ListRules(ctx context.Context) ([]models.ReminderRuleModel, error) {
	return []models.ReminderRuleModel{}, nil
}
func (r *StubReminderRuleRepo) GetRuleByID(ctx context.Context, id uuid.UUID) (*models.ReminderRuleModel, error) {
	return &models.ReminderRuleModel{ID: id, Status: "ACTIVE"}, nil
}
func (r *StubReminderRuleRepo) UpdateRule(ctx context.Context, rule *models.ReminderRuleModel) error { return nil }
func (r *StubReminderRuleRepo) ExecuteReorderTransaction(ctx context.Context, updates map[uuid.UUID]int) error { return nil }

// --- Notification Template Repo ---

type StubNotificationTemplateRepo struct{}

func (r *StubNotificationTemplateRepo) CreateTemplate(ctx context.Context, template *models.NotificationTemplateModel) error { return nil }
func (r *StubNotificationTemplateRepo) ListTemplates(ctx context.Context) ([]models.NotificationTemplateModel, error) {
	return []models.NotificationTemplateModel{}, nil
}
func (r *StubNotificationTemplateRepo) GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.NotificationTemplateModel, error) {
	return &models.NotificationTemplateModel{ID: id, Code: "PAYMENT_REMINDER", Status: "ACTIVE"}, nil
}
func (r *StubNotificationTemplateRepo) UpdateTemplate(ctx context.Context, template *models.NotificationTemplateModel) error { return nil }

// --- Notification Repo & Validator ---

type StubNotificationRepo struct{}

func (r *StubNotificationRepo) CreateNotification(ctx context.Context, notification *models.NotificationModel) error { return nil }
func (r *StubNotificationRepo) GetNotificationByID(ctx context.Context, id uuid.UUID) (*models.NotificationModel, error) {
	return &models.NotificationModel{ID: id, Status: "SENT"}, nil
}
func (r *StubNotificationRepo) GetNotificationByProviderMessageID(ctx context.Context, providerMsgID string) (*models.NotificationModel, error) {
	return nil, ErrNotFound
}
func (r *StubNotificationRepo) ListNotifications(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.NotificationModel, int, error) {
	return []models.NotificationModel{}, 0, nil
}
func (r *StubNotificationRepo) UpdateNotification(ctx context.Context, notification *models.NotificationModel) error { return nil }

type StubNotificationValidator struct{}

func (v *StubNotificationValidator) GetInvoiceOutstandingAmount(ctx context.Context, invoiceID uuid.UUID) (float64, uuid.UUID, error) {
	return 1500.0, uuid.New(), nil
}
func (v *StubNotificationValidator) GetCustomerPhoneNumber(ctx context.Context, customerID uuid.UUID) (string, error) {
	return "+919876543210", nil
}

// --- Report Repo ---

type StubReportRepo struct{}

func (r *StubReportRepo) GetDashboardSummary(ctx context.Context, filters map[string]string) (*models.DashboardReportResponse, error) {
	resp := &models.DashboardReportResponse{}
	resp.Data.GrossSales = "25000.00"
	resp.Data.NetSales = "23000.00"
	resp.Data.CashCollected = "18000.00"
	resp.Data.OutstandingReceivables = "45000.00"
	resp.Data.OverdueAmount = "12000.00"
	resp.Data.COGS = "15000.00"
	resp.Data.GrossProfit = "8000.00"
	resp.Data.StoreExpenses = "3000.00"
	resp.Data.CompanyExpenses = "1000.00"
	resp.Data.NetProfit = "4000.00"
	resp.Data.InventoryValue = "50000.00"
	resp.Data.YellowCustomers = 3
	resp.Data.RedCustomers = 2
	return resp, nil
}
func (r *StubReportRepo) GetSalesReport(ctx context.Context, filters map[string]string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"date": "2026-10-06", "total_sales": 25000}}, nil
}
func (r *StubReportRepo) GetRevenueReport(ctx context.Context, filters map[string]string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"month": "October", "revenue": 150000}}, nil
}
func (r *StubReportRepo) GetCollectionReport(ctx context.Context, filters map[string]string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"date": "2026-10-06", "collected": 18000}}, nil
}
func (r *StubReportRepo) GetExpenseReport(ctx context.Context, filters map[string]string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"category": "Rent", "amount": 10000}}, nil
}
func (r *StubReportRepo) GetProfitReport(ctx context.Context, filters map[string]string) (*models.ProfitReportResponse, error) {
	resp := &models.ProfitReportResponse{}
	resp.Data.Revenue = "150000.00"
	resp.Data.COGS = "90000.00"
	resp.Data.GrossProfit = "60000.00"
	resp.Data.Expenses = "15000.00"
	resp.Data.NetProfit = "45000.00"
	return resp, nil
}
func (r *StubReportRepo) GetReceivablesReport(ctx context.Context, filters map[string]string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"customer": "Acme Ltd", "due": 5000}}, nil
}
func (r *StubReportRepo) GetCreditAgingReport(ctx context.Context, filters map[string]string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"bucket": "0-30 days", "amount": 35000}}, nil
}

// --- Audit Log Repo ---

type StubAuditLogRepo struct{}

func (r *StubAuditLogRepo) ListAuditLogs(ctx context.Context, page, pageSize int, filters map[string]string) ([]models.AuditLogModel, int, error) {
	return []models.AuditLogModel{}, 0, nil
}
func (r *StubAuditLogRepo) GetAuditLogByID(ctx context.Context, id uuid.UUID) (*models.AuditLogModel, error) {
	return &models.AuditLogModel{ID: id, Action: "USER_LOGIN"}, nil
}
func (r *StubAuditLogRepo) GetEntityAuditHistory(ctx context.Context, entityType string, entityID uuid.UUID) ([]models.AuditLogModel, error) {
	return []models.AuditLogModel{}, nil
}
