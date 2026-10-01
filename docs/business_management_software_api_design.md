# API Design — Business Management Software

**Backend:** Golang  
**Database:** PostgreSQL  
**API Style:** REST  
**Base Path:** `/api/v1`  
**Primary Clients:** Web Admin UI, future mobile/admin clients  
**Authentication:** Bearer access token + refresh token  
**Architecture:** Modular Monolith with separate API and background worker processes

---

# 1. API Design Goals

The API should be:

- Consistent across modules.
- Explicit about authorization and store scope.
- Safe for financial operations.
- Idempotent for critical write operations.
- Easy to evolve without breaking clients.
- Compatible with a modular Go backend.
- Clear about business state transitions.
- Friendly for reporting, filtering, and pagination.
- Auditable for all sensitive operations.

Core modules covered:

```text
Authentication
Users
Roles & Permissions
Stores
Customers
Product Categories
Products
Suppliers
Procurement / Purchases
Inventory
Sales / Invoices
Payments
Expenses
Configuration
Reminder Rules
Notification Templates
Notifications
Reporting
Audit Logs
Files / Attachments
System / Health
```

---

# 2. Base URL

```http
/api/v1
```

Example:

```http
GET /api/v1/customers
POST /api/v1/sales
GET /api/v1/reports/dashboard
```

Versioning through the URL keeps future breaking changes manageable.

---

# 3. Authentication

Protected APIs require:

```http
Authorization: Bearer <access_token>
```

Recommended authentication model:

```text
Short-lived access token
+
Long-lived refresh token
```

The backend must enforce permissions and store scope.

Frontend visibility rules are not sufficient authorization.

---

# 4. Common Headers

Recommended headers:

```http
Authorization: Bearer <token>

Content-Type: application/json

X-Request-ID: <uuid>

Idempotency-Key: <uuid>
```

`Idempotency-Key` is required for critical financial writes such as:

```text
Create Sale
Create Payment
Reverse Payment
Approve Purchase
Cancel Sale
Inventory Adjustment
```

---

# 5. Standard Success Response

For single-resource responses:

```json
{
  "data": {
    "id": "uuid"
  }
}
```

For list responses:

```json
{
  "data": [],
  "meta": {
    "page": 1,
    "page_size": 20,
    "total": 145,
    "total_pages": 8
  }
}
```

For actions that do not require a response body:

```http
204 No Content
```

---

# 6. Standard Error Response

All APIs should return a consistent error shape.

```json
{
  "error": {
    "code": "INSUFFICIENT_STOCK",
    "message": "Requested quantity exceeds available stock",
    "details": {
      "product_id": "uuid",
      "requested_quantity": "10.000",
      "available_quantity": "4.000"
    },
    "request_id": "uuid"
  }
}
```

Common HTTP mappings:

| HTTP Status | Meaning |
|---|---|
| 400 | Invalid request |
| 401 | Not authenticated |
| 403 | Not authorized |
| 404 | Resource not found |
| 409 | State conflict / duplicate |
| 422 | Business validation failed |
| 429 | Rate limit exceeded |
| 500 | Internal server error |
| 503 | Dependency temporarily unavailable |

---

# 7. Common Error Codes

```text
VALIDATION_ERROR

UNAUTHORIZED

FORBIDDEN

RESOURCE_NOT_FOUND

DUPLICATE_REQUEST

INVALID_STATE_TRANSITION

CUSTOMER_ALREADY_EXISTS

PRODUCT_ALREADY_EXISTS

INSUFFICIENT_STOCK

PURCHASE_ALREADY_APPROVED

PURCHASE_NOT_APPROVABLE

INVOICE_ALREADY_CANCELLED

INVOICE_ALREADY_PAID

PAYMENT_EXCEEDS_OUTSTANDING

PAYMENT_ALREADY_REVERSED

INVALID_PAYMENT_ALLOCATION

STORE_ACCESS_DENIED

USER_INACTIVE

INVALID_CREDENTIALS

REFRESH_TOKEN_EXPIRED

REMINDER_ALREADY_SENT

FILE_UPLOAD_FAILED

DEPENDENCY_UNAVAILABLE

INTERNAL_ERROR
```

---

# 8. Pagination

Standard query parameters:

```http
?page=1&page_size=20
```

Recommended maximum:

```text
page_size <= 100
```

List responses should include:

```json
{
  "meta": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

---

# 9. Sorting

Recommended format:

```http
?sort_by=created_at&sort_order=desc
```

Allowed order:

```text
asc
desc
```

Each API should explicitly whitelist sortable fields.

---

# 10. Date Filtering

Standardize date filters:

```http
?from_date=2026-10-01&to_date=2026-10-31
```

For timestamps where needed:

```http
?from=2026-10-01T00:00:00+05:30
&to=2026-10-31T23:59:59+05:30
```

---

# 11. Search

Standard search parameter:

```http
?q=rahul
```

Examples:

```http
GET /customers?q=rahul
GET /products?q=asian%20paint
GET /suppliers?q=abc
```

---

# 12. Authentication APIs

## 12.1 Login

```http
POST /api/v1/auth/login
```

Request:

```json
{
  "email": "admin@example.com",
  "password": "secret"
}
```

Response:

```json
{
  "data": {
    "access_token": "jwt",
    "refresh_token": "token",
    "expires_in": 900,
    "user": {
      "id": "uuid",
      "full_name": "Amit",
      "email": "admin@example.com",
      "roles": [
        "STORE_ADMIN"
      ],
      "stores": [
        {
          "id": "uuid",
          "name": "Barpeta Road Store"
        }
      ]
    }
  }
}
```

---

## 12.2 Refresh Access Token

```http
POST /api/v1/auth/refresh
```

Request:

```json
{
  "refresh_token": "token"
}
```

Response:

```json
{
  "data": {
    "access_token": "new-jwt",
    "refresh_token": "rotated-token",
    "expires_in": 900
  }
}
```

Refresh-token rotation is recommended.

---

## 12.3 Logout

```http
POST /api/v1/auth/logout
```

Request:

```json
{
  "refresh_token": "token"
}
```

Response:

```http
204 No Content
```

---

## 12.4 Get Current User

```http
GET /api/v1/auth/me
```

Response:

```json
{
  "data": {
    "id": "uuid",
    "full_name": "Amit",
    "email": "admin@example.com",
    "roles": [
      "STORE_ADMIN"
    ],
    "permissions": [
      "SALE_CREATE",
      "CUSTOMER_CREATE",
      "PAYMENT_CREATE"
    ],
    "stores": [
      {
        "id": "uuid",
        "name": "Store A"
      }
    ]
  }
}
```

---

## 12.5 Change Password

```http
POST /api/v1/auth/change-password
```

Request:

```json
{
  "current_password": "old",
  "new_password": "new"
}
```

---

## 12.6 Forgot Password

```http
POST /api/v1/auth/forgot-password
```

Request:

```json
{
  "email": "admin@example.com"
}
```

Use a generic success response so the endpoint does not reveal whether the account exists.

---

## 12.7 Reset Password

```http
POST /api/v1/auth/reset-password
```

Request:

```json
{
  "token": "reset-token",
  "new_password": "new-password"
}
```

---

# 13. User Management APIs

Primarily Super Admin.

## 13.1 Create User

```http
POST /api/v1/users
```

Request:

```json
{
  "full_name": "Rohan Das",
  "email": "rohan@example.com",
  "phone": "+919876543210",
  "role_codes": [
    "STORE_ADMIN"
  ],
  "store_ids": [
    "store-uuid"
  ]
}
```

Response:

```http
201 Created
```

The system may either:

```text
Generate temporary password
OR
Send account setup/reset link
```

The setup-link approach is preferred.

---

## 13.2 List Users

```http
GET /api/v1/users
```

Filters:

```text
q
status
role
store_id
page
page_size
sort_by
sort_order
```

Example:

```http
GET /api/v1/users?role=STORE_ADMIN&store_id=<uuid>&status=ACTIVE
```

---

## 13.3 Get User

```http
GET /api/v1/users/{user_id}
```

---

## 13.4 Update User

```http
PATCH /api/v1/users/{user_id}
```

Request example:

```json
{
  "full_name": "Rohan D.",
  "phone": "+919999999999"
}
```

---

## 13.5 Activate User

```http
POST /api/v1/users/{user_id}/activate
```

---

## 13.6 Deactivate User

```http
POST /api/v1/users/{user_id}/deactivate
```

Financial records created by the user remain unchanged.

---

## 13.7 Assign Roles

```http
PUT /api/v1/users/{user_id}/roles
```

Request:

```json
{
  "role_codes": [
    "STORE_ADMIN"
  ]
}
```

---

## 13.8 Assign Stores

```http
PUT /api/v1/users/{user_id}/stores
```

Request:

```json
{
  "store_ids": [
    "uuid-1",
    "uuid-2"
  ]
}
```

---

# 14. Role & Permission APIs

Primarily Super Admin.

## 14.1 List Roles

```http
GET /api/v1/roles
```

---

## 14.2 Get Role

```http
GET /api/v1/roles/{role_id}
```

---

## 14.3 List Permissions

```http
GET /api/v1/permissions
```

---

## 14.4 Update Role Permissions

```http
PUT /api/v1/roles/{role_id}/permissions
```

Request:

```json
{
  "permission_codes": [
    "SALE_VIEW",
    "PAYMENT_CREATE",
    "EXPENSE_VIEW"
  ]
}
```

For V1, built-in role creation itself may remain non-configurable.

---

# 15. Store APIs

## 15.1 Create Store

```http
POST /api/v1/stores
```

Permission:

```text
STORE_CREATE
```

Request:

```json
{
  "code": "BRP01",
  "name": "Barpeta Road Main Store",
  "phone": "+919876543210",
  "gstin": "GSTIN",
  "address": {
    "line1": "Main Road",
    "line2": null,
    "city": "Barpeta Road",
    "state": "Assam",
    "postal_code": "781315",
    "country_code": "IN"
  }
}
```

---

## 15.2 List Stores

```http
GET /api/v1/stores
```

Filters:

```text
status
q
```

---

## 15.3 Get Store

```http
GET /api/v1/stores/{store_id}
```

---

## 15.4 Update Store

```http
PATCH /api/v1/stores/{store_id}
```

---

## 15.5 Activate Store

```http
POST /api/v1/stores/{store_id}/activate
```

---

## 15.6 Deactivate Store

```http
POST /api/v1/stores/{store_id}/deactivate
```

A store should not be physically deleted.

---

# 16. Customer APIs

## 16.1 Create Customer

```http
POST /api/v1/customers
```

Request:

```json
{
  "name": "Rahul Das",
  "phone_number": "+919876543210",
  "alternate_phone": null,
  "email": null,
  "address": "Barpeta Road",
  "gstin": null,
  "whatsapp_opt_in": true,
  "store_id": "uuid",
  "notes": null
}
```

Business rules:

```text
phone_number is mandatory
phone_number is company-wide unique
store_id must be accessible to current user
```

If customer exists:

```http
409 CUSTOMER_ALREADY_EXISTS
```

Response should include existing customer ID when useful.

---

## 16.2 List Customers

```http
GET /api/v1/customers
```

Filters:

```text
q
phone
store_id
status
risk_status
has_outstanding
page
page_size
sort_by
sort_order
```

Examples:

```http
GET /api/v1/customers?phone=%2B919876543210
```

```http
GET /api/v1/customers?risk_status=RED&has_outstanding=true
```

---

## 16.3 Get Customer

```http
GET /api/v1/customers/{customer_id}
```

Suggested response:

```json
{
  "data": {
    "id": "uuid",
    "name": "Rahul Das",
    "phone_number": "+919876543210",
    "risk_status": "RED",
    "total_outstanding": "12000.00",
    "oldest_due_date": "2026-09-15",
    "max_overdue_days": 17
  }
}
```

Risk should preferably be calculated from outstanding invoices.

---

## 16.4 Update Customer

```http
PATCH /api/v1/customers/{customer_id}
```

Phone-number change should require careful duplicate validation and audit logging.

---

## 16.5 Deactivate Customer

```http
POST /api/v1/customers/{customer_id}/deactivate
```

Do not delete customer history.

---

## 16.6 Activate Customer

```http
POST /api/v1/customers/{customer_id}/activate
```

---

## 16.7 Customer Store Relationships

```http
GET /api/v1/customers/{customer_id}/stores
```

Optional explicit association API:

```http
POST /api/v1/customers/{customer_id}/stores/{store_id}
```

Usually this association can be created implicitly when the customer is first used by a store.

---

## 16.8 Customer Invoice History

```http
GET /api/v1/customers/{customer_id}/invoices
```

Filters:

```text
payment_status
overdue_only
store_id
from_date
to_date
page
page_size
```

---

## 16.9 Customer Payment History

```http
GET /api/v1/customers/{customer_id}/payments
```

---

## 16.10 Customer Outstanding Summary

```http
GET /api/v1/customers/{customer_id}/outstanding
```

Response:

```json
{
  "data": {
    "customer_id": "uuid",
    "total_outstanding": "12000.00",
    "oldest_due_date": "2026-09-15",
    "max_overdue_days": 17,
    "risk_status": "RED",
    "invoices": [
      {
        "invoice_id": "uuid",
        "invoice_number": "INV-1001",
        "outstanding_amount": "5000.00",
        "due_date": "2026-09-15",
        "overdue_days": 17
      }
    ]
  }
}
```

---

# 17. Product Category APIs

## 17.1 Create Category

```http
POST /api/v1/product-categories
```

Request:

```json
{
  "name": "Paint",
  "description": "Paint products"
}
```

---

## 17.2 List Categories

```http
GET /api/v1/product-categories
```

---

## 17.3 Update Category

```http
PATCH /api/v1/product-categories/{category_id}
```

---

## 17.4 Deactivate Category

```http
POST /api/v1/product-categories/{category_id}/deactivate
```

---

# 18. Product APIs

## 18.1 Create Product

```http
POST /api/v1/products
```

Request:

```json
{
  "sku": "PAINT-ASIAN-20L",
  "name": "Asian Paint 20L",
  "category_id": "uuid",
  "description": null,
  "hsn_code": "3209",
  "unit_of_measure": "BUCKET",
  "default_selling_price": "4500.00",
  "gst_rate": "18.00"
}
```

---

## 18.2 List Products

```http
GET /api/v1/products
```

Filters:

```text
q
sku
category_id
status
page
page_size
```

---

## 18.3 Get Product

```http
GET /api/v1/products/{product_id}
```

---

## 18.4 Update Product

```http
PATCH /api/v1/products/{product_id}
```

---

## 18.5 Activate Product

```http
POST /api/v1/products/{product_id}/activate
```

---

## 18.6 Deactivate Product

```http
POST /api/v1/products/{product_id}/deactivate
```

Products with historical transactions must never be hard-deleted.

---

## 18.7 Product Inventory Across Stores

```http
GET /api/v1/products/{product_id}/inventory
```

Permission:

```text
INVENTORY_VIEW
```

Response:

```json
{
  "data": [
    {
      "store_id": "uuid",
      "store_name": "Store A",
      "available_quantity": "15.000",
      "average_unit_cost": "4200.0000"
    }
  ]
}
```

---

# 19. Supplier APIs

## 19.1 Create Supplier

```http
POST /api/v1/suppliers
```

Request:

```json
{
  "name": "ABC Distributor",
  "phone": "+919876543210",
  "email": "abc@example.com",
  "gstin": "GSTIN",
  "address": "Guwahati"
}
```

---

## 19.2 List Suppliers

```http
GET /api/v1/suppliers
```

Filters:

```text
q
status
page
page_size
```

---

## 19.3 Get Supplier

```http
GET /api/v1/suppliers/{supplier_id}
```

---

## 19.4 Update Supplier

```http
PATCH /api/v1/suppliers/{supplier_id}
```

---

## 19.5 Deactivate Supplier

```http
POST /api/v1/suppliers/{supplier_id}/deactivate
```

---

# 20. Purchase / Procurement APIs

Purchase states:

```text
DRAFT
PENDING_APPROVAL
APPROVED
REJECTED
CANCELLED
```

Allowed transitions:

```text
DRAFT
→ PENDING_APPROVAL
→ APPROVED

DRAFT
→ CANCELLED

PENDING_APPROVAL
→ REJECTED

PENDING_APPROVAL
→ APPROVED
```

A rejected purchase may either remain terminal or be cloned into a new draft.

---

## 20.1 Create Purchase Draft

```http
POST /api/v1/purchases
```

Request:

```json
{
  "store_id": "uuid",
  "supplier_id": "uuid",
  "supplier_invoice_number": "SUP-123",
  "purchase_date": "2026-10-02",
  "items": [
    {
      "product_id": "uuid",
      "quantity": "100.000",
      "base_unit_cost": "100.0000",
      "gst_rate": "18.00"
    },
    {
      "product_id": "uuid-2",
      "quantity": "50.000",
      "base_unit_cost": "200.0000",
      "gst_rate": "18.00"
    }
  ],
  "additional_costs": [
    {
      "cost_type": "TRANSPORTATION",
      "description": "Truck from Guwahati",
      "amount": "5000.00"
    }
  ],
  "notes": null
}
```

The backend should calculate:

```text
base_amount
gst_amount
additional_cost_total
landed_unit_cost
grand_total
```

Do not trust client-computed totals.

---

## 20.2 List Purchases

```http
GET /api/v1/purchases
```

Filters:

```text
store_id
supplier_id
status
from_date
to_date
created_by
page
page_size
```

---

## 20.3 Get Purchase

```http
GET /api/v1/purchases/{purchase_id}
```

---

## 20.4 Update Purchase Draft

```http
PUT /api/v1/purchases/{purchase_id}
```

Only allowed while:

```text
status = DRAFT
```

---

## 20.5 Submit Purchase for Approval

```http
POST /api/v1/purchases/{purchase_id}/submit
```

Permission:

```text
PURCHASE_CREATE
```

State:

```text
DRAFT → PENDING_APPROVAL
```

Creates approval-history event:

```text
SUBMITTED
```

---

## 20.6 Approve Purchase

```http
POST /api/v1/purchases/{purchase_id}/approve
```

Permission:

```text
PURCHASE_APPROVE
```

Header:

```http
Idempotency-Key: <uuid>
```

Business transaction:

```text
Lock purchase

Verify PENDING_APPROVAL

Calculate landed costs

Lock inventory rows

Update weighted average cost

Increase inventory

Create inventory movements

Mark purchase APPROVED

Create approval history

Create audit log
```

Store Admin must not approve their own purchase.

---

## 20.7 Reject Purchase

```http
POST /api/v1/purchases/{purchase_id}/reject
```

Request:

```json
{
  "reason": "Incorrect purchase quantity"
}
```

State:

```text
PENDING_APPROVAL → REJECTED
```

Inventory remains unchanged.

---

## 20.8 Cancel Purchase Draft

```http
POST /api/v1/purchases/{purchase_id}/cancel
```

Allowed only when business rules permit.

Typically:

```text
DRAFT → CANCELLED
```

---

## 20.9 Purchase Approval History

```http
GET /api/v1/purchases/{purchase_id}/history
```

Response includes:

```text
SUBMITTED
APPROVED
REJECTED
CANCELLED
```

with actor and timestamp.

---

# 21. Inventory APIs

## 21.1 Get Store Inventory

```http
GET /api/v1/stores/{store_id}/inventory
```

Filters:

```text
q
product_id
category_id
low_stock
page
page_size
sort_by
sort_order
```

Response item:

```json
{
  "product_id": "uuid",
  "sku": "PAINT-ASIAN-20L",
  "product_name": "Asian Paint 20L",
  "available_quantity": "20.000",
  "average_unit_cost": "4200.5000",
  "inventory_value": "84010.00"
}
```

---

## 21.2 Get Product Inventory in Store

```http
GET /api/v1/stores/{store_id}/inventory/{product_id}
```

---

## 21.3 Inventory Movement History

```http
GET /api/v1/stores/{store_id}/inventory/{product_id}/movements
```

Filters:

```text
movement_type
from
to
page
page_size
```

---

## 21.4 Create Inventory Adjustment

```http
POST /api/v1/stores/{store_id}/inventory/adjustments
```

Permission:

```text
INVENTORY_ADJUST
```

Header:

```http
Idempotency-Key: <uuid>
```

Request:

```json
{
  "product_id": "uuid",
  "adjustment_type": "ADJUSTMENT_OUT",
  "quantity": "3.000",
  "reason": "Physical stock damaged"
}
```

The backend should:

```text
Lock inventory
Validate resulting quantity >= 0
Update inventory
Create inventory movement
Create audit log
```

Do not expose direct arbitrary inventory update APIs.

---

## 21.5 Low Stock Report

```http
GET /api/v1/stores/{store_id}/inventory/low-stock
```

For V1, low-stock threshold may be a query parameter or product configuration.

Example:

```http
GET /api/v1/stores/{store_id}/inventory/low-stock?threshold=5
```

---

# 22. Sales / Invoice APIs

Invoice states:

```text
CONFIRMED
CANCELLED
```

Payment states:

```text
UNPAID
PARTIALLY_PAID
PAID
```

---

## 22.1 Create Sale

```http
POST /api/v1/sales
```

Header:

```http
Idempotency-Key: <uuid>
```

Request:

```json
{
  "store_id": "uuid",
  "customer_id": "uuid",
  "invoice_date": "2026-10-02",
  "items": [
    {
      "product_id": "uuid",
      "quantity": "2.000",
      "unit_price": "1000.00",
      "discount_amount": "50.00"
    },
    {
      "product_id": "uuid-2",
      "quantity": "1.000",
      "unit_price": "500.00",
      "discount_amount": "0.00"
    }
  ],
  "initial_payment": {
    "amount": "1000.00",
    "payment_method": "UPI",
    "transaction_reference": "TXN123"
  },
  "due_date": "2026-10-17",
  "notes": null
}
```

Backend responsibilities:

```text
Validate store access

Validate customer

Sort product IDs

Lock inventory rows

Validate stock

Read current average_unit_cost

Calculate invoice values

Create invoice

Create invoice items

Reduce inventory

Create inventory movements

Create initial payment if amount > 0

Create payment allocation

Set paid/outstanding amounts

Create audit record
```

Client should not provide:

```text
subtotal
GST total
grand total
COGS
paid_amount
outstanding_amount
payment_status
```

The server calculates these.

---

## 22.2 Create Cash Sale Without Customer

This is a product decision.

For the current PRD, all sales appear customer-based.

Recommended V1:

```text
customer_id is mandatory
```

If anonymous cash sales are needed later, either:

```text
Allow customer_id nullable
```

or create a store-level:

```text
WALK_IN_CUSTOMER
```

Do not silently introduce this until business confirms it.

---

## 22.3 List Sales

```http
GET /api/v1/sales
```

Filters:

```text
store_id
customer_id
payment_status
status
overdue_only
from_date
to_date
created_by
page
page_size
sort_by
sort_order
```

Example:

```http
GET /api/v1/sales?store_id=<uuid>&payment_status=PARTIALLY_PAID
```

---

## 22.4 Get Sale

```http
GET /api/v1/sales/{sale_id}
```

Response should include:

```text
Invoice header

Invoice items

Customer

Store

Payment summary

Outstanding amount

Due date

Overdue days

Created by

Cancellation info if any
```

---

## 22.5 Get Invoice by Invoice Number

Optional convenience API:

```http
GET /api/v1/sales/by-number/{invoice_number}
```

---

## 22.6 Cancel Sale

```http
POST /api/v1/sales/{sale_id}/cancel
```

Header:

```http
Idempotency-Key: <uuid>
```

Request:

```json
{
  "reason": "Incorrect invoice"
}
```

Cancellation should be tightly controlled.

The backend must consider:

```text
Inventory restoration

Existing payments

Payment reversal / refund requirements

Audit history
```

Recommended V1 rule:

A sale cannot be cancelled automatically when there are unreversed payments unless the payment reversal/refund workflow is completed first.

Cancellation creates reversing inventory movements rather than deleting the original movements.

---

## 22.7 Sale Payment Summary

```http
GET /api/v1/sales/{sale_id}/payments
```

Returns allocations applied to the invoice.

---

## 22.8 Sale Inventory Movements

```http
GET /api/v1/sales/{sale_id}/inventory-movements
```

Useful for audit/debugging.

---

# 23. Payment APIs

## 23.1 Collect Customer Payment

```http
POST /api/v1/payments
```

Header:

```http
Idempotency-Key: <uuid>
```

Request with automatic oldest-first allocation:

```json
{
  "customer_id": "uuid",
  "amount": "6000.00",
  "payment_method": "UPI",
  "transaction_reference": "TXN123",
  "collection_store_id": null,
  "allocation_strategy": "OLDEST_FIRST",
  "notes": null
}
```

Request with explicit allocation:

```json
{
  "customer_id": "uuid",
  "amount": "6000.00",
  "payment_method": "BANK_TRANSFER",
  "transaction_reference": "BANK-REF",
  "allocation_strategy": "EXPLICIT",
  "allocations": [
    {
      "invoice_id": "invoice-1",
      "amount": "5000.00"
    },
    {
      "invoice_id": "invoice-2",
      "amount": "1000.00"
    }
  ]
}
```

Backend transaction:

```text
Lock outstanding invoices

Validate invoice belongs to customer

Create payment

Create allocations

Update invoice paid_amount

Update invoice outstanding_amount

Update payment_status

Create audit log
```

---

## 23.2 List Payments

```http
GET /api/v1/payments
```

Filters:

```text
customer_id
invoice_id
collection_store_id
collected_by
payment_method
status
from
to
page
page_size
```

---

## 23.3 Get Payment

```http
GET /api/v1/payments/{payment_id}
```

Should include allocation details.

---

## 23.4 Reverse Payment

```http
POST /api/v1/payments/{payment_id}/reverse
```

Header:

```http
Idempotency-Key: <uuid>
```

Request:

```json
{
  "reason": "Payment entered incorrectly"
}
```

Backend should:

```text
Lock payment

Ensure POSTED

Lock affected invoices

Reverse allocations

Increase invoice outstanding

Update payment statuses

Create reversal record

Create audit log
```

Never delete the original payment.

---

## 23.5 Get Payment Allocations

```http
GET /api/v1/payments/{payment_id}/allocations
```

---

# 24. Expense Category APIs

## 24.1 Create Expense Category

```http
POST /api/v1/expense-categories
```

Request:

```json
{
  "name": "Electricity",
  "scope": "STORE"
}
```

---

## 24.2 List Expense Categories

```http
GET /api/v1/expense-categories
```

Filters:

```text
scope
status
```

---

## 24.3 Update Expense Category

```http
PATCH /api/v1/expense-categories/{category_id}
```

---

## 24.4 Deactivate Expense Category

```http
POST /api/v1/expense-categories/{category_id}/deactivate
```

---

# 25. Expense APIs

## 25.1 Create Expense

```http
POST /api/v1/expenses
```

Request:

```json
{
  "store_id": "uuid",
  "category_id": "uuid",
  "expense_date": "2026-10-02",
  "amount": "3000.00",
  "gst_amount": "0.00",
  "payment_method": "UPI",
  "vendor": "Electricity Board",
  "description": "September electricity bill",
  "receipt_object_key": null
}
```

For company expense:

```json
{
  "store_id": null,
  "category_id": "uuid",
  "expense_date": "2026-10-02",
  "amount": "10000.00",
  "description": "CA fees"
}
```

Only authorized roles may create company-level expenses.

---

## 25.2 List Expenses

```http
GET /api/v1/expenses
```

Filters:

```text
store_id
category_id
scope
status
created_by
from_date
to_date
page
page_size
```

---

## 25.3 Get Expense

```http
GET /api/v1/expenses/{expense_id}
```

---

## 25.4 Update Expense

```http
PATCH /api/v1/expenses/{expense_id}
```

Recommended rule:

Allow modification only while the expense is operationally editable.

For finalized financial records, prefer reversal instead.

---

## 25.5 Reverse Expense

```http
POST /api/v1/expenses/{expense_id}/reverse
```

Request:

```json
{
  "reason": "Duplicate expense entry"
}
```

Original expense remains in history.

---

# 26. Company Configuration APIs

Primarily Super Admin.

## 26.1 Get Company Settings

```http
GET /api/v1/configuration/company
```

Response:

```json
{
  "data": {
    "company_name": "Example Business",
    "currency": "INR",
    "timezone": "Asia/Kolkata",
    "default_credit_days": 15,
    "yellow_after_overdue_days": 1,
    "red_after_overdue_days": 7,
    "additional_cost_allocation_method": "PURCHASE_VALUE",
    "include_gst_in_landed_cost": false
  }
}
```

---

## 26.2 Update Company Settings

```http
PATCH /api/v1/configuration/company
```

Request:

```json
{
  "default_credit_days": 15,
  "yellow_after_overdue_days": 1,
  "red_after_overdue_days": 7
}
```

Validation:

```text
yellow >= 0

red > yellow
```

All changes should be audited.

---

# 27. Reminder Rule APIs

## 27.1 Create Reminder Rule

```http
POST /api/v1/reminder-rules
```

Request:

```json
{
  "sequence_number": 1,
  "days_after_due": 7,
  "template_id": "uuid",
  "status": "ACTIVE"
}
```

---

## 27.2 List Reminder Rules

```http
GET /api/v1/reminder-rules
```

---

## 27.3 Update Reminder Rule

```http
PATCH /api/v1/reminder-rules/{rule_id}
```

---

## 27.4 Activate Reminder Rule

```http
POST /api/v1/reminder-rules/{rule_id}/activate
```

---

## 27.5 Deactivate Reminder Rule

```http
POST /api/v1/reminder-rules/{rule_id}/deactivate
```

---

## 27.6 Reorder Reminder Rules

Optional:

```http
PUT /api/v1/reminder-rules/order
```

Request:

```json
{
  "rules": [
    {
      "id": "uuid-1",
      "sequence_number": 1
    },
    {
      "id": "uuid-2",
      "sequence_number": 2
    }
  ]
}
```

---

# 28. Notification Template APIs

## 28.1 Create Notification Template

```http
POST /api/v1/notification-templates
```

Request:

```json
{
  "code": "PAYMENT_DUE_REMINDER",
  "provider_template_id": "provider-template",
  "language_code": "en",
  "body_preview": "Your payment is overdue...",
  "status": "ACTIVE"
}
```

---

## 28.2 List Notification Templates

```http
GET /api/v1/notification-templates
```

---

## 28.3 Update Notification Template

```http
PATCH /api/v1/notification-templates/{template_id}
```

---

# 29. Notification APIs

## 29.1 Send Manual Reminder

```http
POST /api/v1/invoices/{invoice_id}/reminders
```

Permission:

```text
REMINDER_SEND
```

Request:

```json
{
  "template_id": "uuid",
  "notes": null
}
```

Backend should:

```text
Validate invoice outstanding > 0

Validate WhatsApp consent / permitted messaging policy

Create notification record

Create outbox event or enqueue worker work

Return accepted response
```

Response:

```http
202 Accepted
```

Example:

```json
{
  "data": {
    "notification_id": "uuid",
    "status": "PENDING"
  }
}
```

---

## 29.2 List Notifications

```http
GET /api/v1/notifications
```

Filters:

```text
customer_id
invoice_id
trigger_type
status
notification_type
from
to
page
page_size
```

---

## 29.3 Get Notification

```http
GET /api/v1/notifications/{notification_id}
```

---

## 29.4 Retry Failed Notification

```http
POST /api/v1/notifications/{notification_id}/retry
```

Only valid for retryable failed notifications.

---

## 29.5 Customer Notification History

```http
GET /api/v1/customers/{customer_id}/notifications
```

---

## 29.6 WhatsApp Provider Webhook

```http
POST /api/v1/webhooks/whatsapp
```

Used for provider delivery status callbacks.

Must not use normal admin authentication.

Instead validate provider-specific:

```text
Signature
Secret
Webhook verification token
```

Possible events:

```text
SENT
DELIVERED
FAILED
READ
```

Update notification state idempotently.

---

# 30. Reporting APIs

Reporting APIs should be read-only and support store and date filters.

---

## 30.1 Dashboard Summary

```http
GET /api/v1/reports/dashboard
```

Filters:

```text
store_id
from_date
to_date
```

Response:

```json
{
  "data": {
    "gross_sales": "500000.00",
    "net_sales": "480000.00",
    "cash_collected": "420000.00",
    "outstanding_receivables": "60000.00",
    "overdue_amount": "25000.00",
    "cogs": "300000.00",
    "gross_profit": "180000.00",
    "store_expenses": "40000.00",
    "company_expenses": "20000.00",
    "net_profit": "120000.00",
    "inventory_value": "850000.00",
    "yellow_customers": 12,
    "red_customers": 5
  }
}
```

---

## 30.2 Sales Report

```http
GET /api/v1/reports/sales
```

Filters:

```text
store_id
customer_id
from_date
to_date
group_by
```

Suggested `group_by`:

```text
DAY
MONTH
STORE
PRODUCT
CUSTOMER
```

---

## 30.3 Revenue Report

```http
GET /api/v1/reports/revenue
```

---

## 30.4 Collection Report

```http
GET /api/v1/reports/collections
```

This must remain separate from revenue.

Filters:

```text
store_id
collection_store_id
payment_method
from_date
to_date
group_by
```

---

## 30.5 Expense Report

```http
GET /api/v1/reports/expenses
```

Filters:

```text
store_id
category_id
scope
from_date
to_date
group_by
```

---

## 30.6 Profit Report

```http
GET /api/v1/reports/profit
```

Response concept:

```json
{
  "data": {
    "revenue": "480000.00",
    "cogs": "300000.00",
    "gross_profit": "180000.00",
    "expenses": "60000.00",
    "net_profit": "120000.00"
  }
}
```

---

## 30.7 Receivables Report

```http
GET /api/v1/reports/receivables
```

Filters:

```text
store_id
risk_status
overdue_only
customer_id
```

---

## 30.8 Credit Aging Report

```http
GET /api/v1/reports/receivables/aging
```

Response:

```json
{
  "data": {
    "buckets": [
      {
        "label": "NOT_DUE",
        "amount": "200000.00"
      },
      {
        "label": "1_7_DAYS",
        "amount": "50000.00"
      },
      {
        "label": "8_30_DAYS",
        "amount": "30000.00"
      },
      {
        "label": "31_60_DAYS",
        "amount": "15000.00"
      },
      {
        "label": "60_PLUS_DAYS",
        "amount": "10000.00"
      }
    ]
  }
}
```

---

## 30.9 Store Comparison Report

```http
GET /api/v1/reports/stores
```

Response per store may include:

```text
Sales
Collections
Outstanding
COGS
Gross Profit
Expenses
Net Profit
Inventory Value
```

---

## 30.10 Product Sales Report

```http
GET /api/v1/reports/products
```

Filters:

```text
store_id
product_id
category_id
from_date
to_date
```

Metrics:

```text
quantity_sold
sales_value
cogs
gross_profit
```

---

## 30.11 Inventory Valuation Report

```http
GET /api/v1/reports/inventory-value
```

Response by store/product:

```text
available_quantity
average_unit_cost
inventory_value
```

---

# 31. Audit Log APIs

## 31.1 List Audit Logs

```http
GET /api/v1/audit-logs
```

Permissions should be restricted.

Filters:

```text
actor_user_id
store_id
action
entity_type
entity_id
from
to
page
page_size
```

---

## 31.2 Get Audit Log

```http
GET /api/v1/audit-logs/{audit_log_id}
```

---

## 31.3 Entity Audit History

Optional convenience endpoint:

```http
GET /api/v1/audit-logs/entity/{entity_type}/{entity_id}
```

Example:

```http
GET /api/v1/audit-logs/entity/PAYMENT/<payment-id>
```

---

# 32. File / Attachment APIs

Used for:

```text
Expense receipts
Purchase invoices
Future documents
```

Recommended approach:

```text
Client
→ request upload authorization
→ upload to object storage
→ store object key in domain entity
```

---

## 32.1 Request Upload URL

```http
POST /api/v1/files/upload-url
```

Request:

```json
{
  "file_name": "electricity-bill.pdf",
  "content_type": "application/pdf",
  "purpose": "EXPENSE_RECEIPT"
}
```

Response:

```json
{
  "data": {
    "upload_url": "signed-url",
    "object_key": "expenses/uuid/file.pdf",
    "expires_in": 300
  }
}
```

---

## 32.2 Request Download URL

```http
GET /api/v1/files/download-url?object_key=<encoded-key>
```

Authorization must verify that the caller has access to the underlying record.

---

# 33. System APIs

## 33.1 Liveness

```http
GET /health/live
```

Response:

```http
200 OK
```

Should not depend heavily on external systems.

---

## 33.2 Readiness

```http
GET /health/ready
```

Check at least:

```text
PostgreSQL connectivity
Required configuration loaded
```

Optional dependency health can be exposed carefully.

---

## 33.3 Version

```http
GET /api/v1/system/version
```

Response:

```json
{
  "data": {
    "version": "1.0.0",
    "commit": "abc123"
  }
}
```

---

# 34. Authorization Matrix

A simplified default matrix:

| Module / Action | Store Admin | Finance Admin | Super Admin |
|---|---:|---:|---:|
| Login | Yes | Yes | Yes |
| View own store | Yes | All stores | All stores |
| Create customer | Yes | Optional | Yes |
| View customers | Own scope | All | All |
| Create product | Usually No | No | Yes |
| Create purchase | Yes | No | Yes |
| Approve purchase | No | No | Yes |
| View inventory | Own store | Optional read | All |
| Adjust inventory | Permission-based | No | Yes |
| Create sale | Own store | No | Yes |
| View sales | Own store | All | All |
| Collect payment | Own store | All | All |
| Reverse payment | Usually No | Permission-based | Yes |
| Create store expense | Yes | Optional | Yes |
| Create company expense | No | Optional | Yes |
| Send reminder | Optional | Yes | Yes |
| Configure reminder rules | No | No | Yes |
| View financial reports | Limited | All | All |
| Manage users | No | No | Yes |
| View audit logs | Limited/No | Limited | Yes |

This should still be implemented with permissions rather than hard-coded role checks.

---

# 35. Store Scope Enforcement

For store-scoped endpoints:

```http
GET /api/v1/stores/{store_id}/inventory
```

backend logic should verify:

```text
User has required permission
AND
User has access to store
```

except for globally scoped roles.

Do not accept a `store_id` merely because the frontend sent it.

---

# 36. Idempotency Design

Critical endpoints:

```text
POST /sales

POST /payments

POST /payments/{id}/reverse

POST /purchases/{id}/approve

POST /sales/{id}/cancel

POST /stores/{store_id}/inventory/adjustments
```

should accept:

```http
Idempotency-Key
```

Rules:

```text
Same user + same key + same request
→ return original response

Same user + same key + different payload
→ 409 DUPLICATE_REQUEST
```

Store:

```text
key
user
endpoint
request_hash
response
expiry
```

---

# 37. API Transaction Boundaries

## Create Sale

One database transaction:

```text
Lock inventory

Create invoice

Create invoice items

Update inventory

Create inventory movements

Create initial payment

Create allocation

Update invoice balance

Create audit entry

Commit
```

---

## Approve Purchase

```text
Lock purchase

Validate state

Calculate landed cost

Lock inventory

Update weighted average

Create inventory movements

Update purchase

Create approval history

Create audit entry

Commit
```

---

## Collect Payment

```text
Lock invoices

Create payment

Create allocations

Update invoice balances

Create audit entry

Commit
```

---

## Reverse Payment

```text
Lock original payment

Lock invoices

Create reversal

Reverse allocations

Restore invoice outstanding

Create audit entry

Commit
```

---

# 38. API State Transition Rules

## Purchase

```text
DRAFT
  ↓
PENDING_APPROVAL
  ↓
APPROVED
```

Alternative:

```text
PENDING_APPROVAL
  ↓
REJECTED
```

---

## Invoice

```text
CONFIRMED
  ↓
CANCELLED
```

No direct deletion.

---

## Payment

```text
POSTED
  ↓
REVERSED
```

No deletion.

---

## Notification

```text
PENDING
   ↓
PROCESSING
   ↓
SENT
   ↓
DELIVERED
```

Failure:

```text
PROCESSING
   ↓
FAILED
   ↓
retry
```

---

# 39. Derived Fields

Some values should be server-derived, not client-controlled.

## Invoice

```text
subtotal
discount total
GST
total_amount
paid_amount
outstanding_amount
payment_status
COGS
overdue status
overdue days
```

## Purchase

```text
base totals
GST totals
allocated additional cost
landed unit cost
grand total
```

## Customer

```text
total outstanding
oldest due date
max overdue days
risk status
```

## Reports

```text
revenue
cash collected
COGS
gross profit
net profit
inventory value
```

---

# 40. Monetary Representation in APIs

Avoid JSON numeric floating-point ambiguity for monetary values.

Recommended:

```json
{
  "amount": "1250.50"
}
```

instead of:

```json
{
  "amount": 1250.5
}
```

Use decimal strings for:

```text
money
unit cost
quantity
tax rates if precision matters
```

Example:

```json
{
  "quantity": "2.500",
  "unit_price": "105.50",
  "gst_rate": "18.00"
}
```

---

# 41. Date and Time Representation

Business dates:

```text
YYYY-MM-DD
```

Example:

```text
2026-10-02
```

Timestamps:

```text
RFC3339 / ISO 8601 with timezone
```

Example:

```text
2026-10-02T16:30:00+05:30
```

The company timezone should come from configuration.

---

# 42. Validation Rules

Examples:

```text
Phone must be normalized.

Email must be valid.

Quantity > 0.

Amount > 0.

GST rate between 0 and 100.

Due date required if outstanding > 0.

Red-risk threshold > yellow-risk threshold.

Inventory adjustment must include a reason.

Cancellation/reversal must include a reason.

Store Admin cannot approve their own purchase.

Payment allocations cannot exceed payment amount.

Payment allocation cannot exceed invoice outstanding.

Sale quantity cannot exceed current available inventory.
```

---

# 43. API Logging

Every request should log:

```text
request_id
user_id
store_id if applicable
method
route
status
duration
```

Sensitive values must not be logged:

```text
password
refresh token
access token
full payment secrets
```

For critical mutations, log domain identifiers:

```text
invoice_id
payment_id
purchase_id
customer_id
```

---

# 44. Audit Events

Important actions should generate audit records.

Examples:

```text
USER_CREATED

USER_DEACTIVATED

STORE_CREATED

CUSTOMER_CREATED

CUSTOMER_UPDATED

PRODUCT_CREATED

PURCHASE_CREATED

PURCHASE_SUBMITTED

PURCHASE_APPROVED

PURCHASE_REJECTED

INVENTORY_ADJUSTED

SALE_CREATED

SALE_CANCELLED

PAYMENT_CREATED

PAYMENT_REVERSED

EXPENSE_CREATED

EXPENSE_REVERSED

REMINDER_SENT_MANUALLY

COMPANY_SETTINGS_UPDATED
```

---

# 45. Background Worker Responsibilities

The following are not normal admin APIs, but are part of the API/system design:

```text
Overdue invoice scan

Reminder-rule evaluation

Notification creation

WhatsApp sending

Retry handling

Outbox processing

Delivery webhook processing
```

The worker and API share the same PostgreSQL database.

---

# 46. Outbox Events

Suggested event types:

```text
PURCHASE_APPROVED

SALE_CREATED

SALE_CANCELLED

PAYMENT_COLLECTED

PAYMENT_REVERSED

INVOICE_OVERDUE

REMINDER_REQUESTED

NOTIFICATION_CREATED
```

The API transaction should insert the outbox event in the same transaction as the business data when atomicity is required.

---

# 47. Recommended Go Handler Structure

Example:

```text
internal/

  customers/
    handler.go
    service.go
    repository.go
    model.go

  sales/
    handler.go
    service.go
    repository.go
    model.go
```

HTTP layer:

```text
Parse request
Validate input
Authenticate
Authorize
Call use case/service
Map domain error → HTTP response
```

The HTTP handler should not own business transactions.

---

# 48. Suggested Route Registration

Conceptually:

```text
/api/v1/auth/*

/api/v1/users/*
/api/v1/roles/*
/api/v1/permissions/*

/api/v1/stores/*

/api/v1/customers/*

/api/v1/product-categories/*
/api/v1/products/*
/api/v1/suppliers/*

/api/v1/purchases/*

/api/v1/stores/{store_id}/inventory/*

/api/v1/sales/*

/api/v1/payments/*

/api/v1/expense-categories/*
/api/v1/expenses/*

/api/v1/configuration/*
/api/v1/reminder-rules/*
/api/v1/notification-templates/*
/api/v1/notifications/*

/api/v1/reports/*

/api/v1/audit-logs/*

/api/v1/files/*

/api/v1/webhooks/*

/health/*
```

---

# 49. Complete Endpoint Summary

## Authentication

```http
POST /auth/login
POST /auth/refresh
POST /auth/logout
GET  /auth/me
POST /auth/change-password
POST /auth/forgot-password
POST /auth/reset-password
```

## Users

```http
POST  /users
GET   /users
GET   /users/{id}
PATCH /users/{id}
POST  /users/{id}/activate
POST  /users/{id}/deactivate
PUT   /users/{id}/roles
PUT   /users/{id}/stores
```

## Roles & Permissions

```http
GET /roles
GET /roles/{id}
GET /permissions
PUT /roles/{id}/permissions
```

## Stores

```http
POST  /stores
GET   /stores
GET   /stores/{id}
PATCH /stores/{id}
POST  /stores/{id}/activate
POST  /stores/{id}/deactivate
```

## Customers

```http
POST  /customers
GET   /customers
GET   /customers/{id}
PATCH /customers/{id}
POST  /customers/{id}/activate
POST  /customers/{id}/deactivate
GET   /customers/{id}/stores
POST  /customers/{id}/stores/{store_id}
GET   /customers/{id}/invoices
GET   /customers/{id}/payments
GET   /customers/{id}/outstanding
GET   /customers/{id}/notifications
```

## Product Categories

```http
POST  /product-categories
GET   /product-categories
PATCH /product-categories/{id}
POST  /product-categories/{id}/deactivate
```

## Products

```http
POST  /products
GET   /products
GET   /products/{id}
PATCH /products/{id}
POST  /products/{id}/activate
POST  /products/{id}/deactivate
GET   /products/{id}/inventory
```

## Suppliers

```http
POST  /suppliers
GET   /suppliers
GET   /suppliers/{id}
PATCH /suppliers/{id}
POST  /suppliers/{id}/deactivate
```

## Purchases

```http
POST /purchases
GET  /purchases
GET  /purchases/{id}
PUT  /purchases/{id}
POST /purchases/{id}/submit
POST /purchases/{id}/approve
POST /purchases/{id}/reject
POST /purchases/{id}/cancel
GET  /purchases/{id}/history
```

## Inventory

```http
GET  /stores/{store_id}/inventory
GET  /stores/{store_id}/inventory/{product_id}
GET  /stores/{store_id}/inventory/{product_id}/movements
POST /stores/{store_id}/inventory/adjustments
GET  /stores/{store_id}/inventory/low-stock
```

## Sales

```http
POST /sales
GET  /sales
GET  /sales/{id}
GET  /sales/by-number/{invoice_number}
POST /sales/{id}/cancel
GET  /sales/{id}/payments
GET  /sales/{id}/inventory-movements
```

## Payments

```http
POST /payments
GET  /payments
GET  /payments/{id}
POST /payments/{id}/reverse
GET  /payments/{id}/allocations
```

## Expense Categories

```http
POST  /expense-categories
GET   /expense-categories
PATCH /expense-categories/{id}
POST  /expense-categories/{id}/deactivate
```

## Expenses

```http
POST  /expenses
GET   /expenses
GET   /expenses/{id}
PATCH /expenses/{id}
POST  /expenses/{id}/reverse
```

## Company Configuration

```http
GET   /configuration/company
PATCH /configuration/company
```

## Reminder Rules

```http
POST  /reminder-rules
GET   /reminder-rules
PATCH /reminder-rules/{id}
POST  /reminder-rules/{id}/activate
POST  /reminder-rules/{id}/deactivate
PUT   /reminder-rules/order
```

## Notification Templates

```http
POST  /notification-templates
GET   /notification-templates
PATCH /notification-templates/{id}
```

## Notifications

```http
POST /invoices/{invoice_id}/reminders
GET  /notifications
GET  /notifications/{id}
POST /notifications/{id}/retry
GET  /customers/{customer_id}/notifications
```

## Webhooks

```http
POST /webhooks/whatsapp
```

## Reports

```http
GET /reports/dashboard
GET /reports/sales
GET /reports/revenue
GET /reports/collections
GET /reports/expenses
GET /reports/profit
GET /reports/receivables
GET /reports/receivables/aging
GET /reports/stores
GET /reports/products
GET /reports/inventory-value
```

## Audit

```http
GET /audit-logs
GET /audit-logs/{id}
GET /audit-logs/entity/{entity_type}/{entity_id}
```

## Files

```http
POST /files/upload-url
GET  /files/download-url
```

## System

```http
GET /health/live
GET /health/ready
GET /api/v1/system/version
```

---

# 50. Recommended Implementation Order

A practical API implementation sequence is:

```text
1. Authentication

2. Users / Roles / Stores

3. Customers

4. Products / Categories / Suppliers

5. Procurement

6. Inventory

7. Sales

8. Payments

9. Expenses

10. Configuration

11. Notification Templates / Reminder Rules

12. Notifications / WhatsApp

13. Reporting

14. Audit APIs

15. Files

16. Operational APIs
```

This order follows domain dependencies and minimizes rework.

---

# 51. Most Critical APIs

The most sensitive endpoints in the system are:

```text
POST /sales

POST /payments

POST /payments/{id}/reverse

POST /purchases/{id}/approve

POST /sales/{id}/cancel

POST /stores/{store_id}/inventory/adjustments
```

These require the strongest guarantees around:

```text
Authorization

Idempotency

PostgreSQL transactions

Row locking

Validation

Audit logging
```

If these APIs are correct, most of the financial and inventory integrity of the system remains protected.

---

# 52. Final API Design Principles

The system should follow these rules consistently:

```text
Server calculates financial totals.

Client never directly updates inventory quantity.

Client never directly updates outstanding amount.

Financial records are reversed, not deleted.

Critical writes are idempotent.

Store scope is always enforced server-side.

State transitions happen through explicit action APIs.

Derived financial values are read-only.

Async notification failures never roll back financial transactions.

All important mutations are audited.

Money is represented as decimal strings in JSON.

Dates and timestamps use standardized formats.

List APIs support pagination and filtering.

Reporting remains separate from transaction APIs.
```

This API surface is intentionally designed to map cleanly onto the earlier PRD, HLD, and PostgreSQL schema.
