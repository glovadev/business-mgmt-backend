# High-Level Design — Business Management Software

**Backend:** Golang  
**Database:** PostgreSQL  
**Recommended architecture:** Modular Monolith  
**Primary external integration:** WhatsApp messaging provider  
**Deployment model:** Stateless API + background worker(s) + PostgreSQL

---

## 1. Architecture Decision

For this product, I would **not start with microservices**.

The system has several domains:

```text
Customer
Inventory
Procurement
Sales
Credit
Payments
Expenses
Notifications
```

but many important operations require strong consistency across those domains.

For example:

```text
Create Sale
   ↓
Validate Inventory
   ↓
Deduct Inventory
   ↓
Create Invoice
   ↓
Record Initial Payment
   ↓
Create Customer Credit
```

Splitting Sales, Inventory, and Payments into separate services immediately would introduce distributed transactions, message brokers, eventual consistency, retry logic, and much higher operational complexity.

Instead:

```text
                    ┌─────────────────────┐
                    │      Web Client     │
                    │ Admin / Finance UI  │
                    └──────────┬──────────┘
                               │ HTTPS
                               ▼
                    ┌─────────────────────┐
                    │    Load Balancer    │
                    │      / Nginx        │
                    └──────────┬──────────┘
                               │
                               ▼
              ┌────────────────────────────────┐
              │           Go Backend           │
              │                                │
              │  ┌──────────────────────────┐  │
              │  │ Authentication / RBAC    │  │
              │  │ Store                    │  │
              │  │ Customer                 │  │
              │  │ Product                  │  │
              │  │ Procurement              │  │
              │  │ Inventory                │  │
              │  │ Sales                    │  │
              │  │ Credit                   │  │
              │  │ Payment                  │  │
              │  │ Expenses                 │  │
              │  │ Notification             │  │
              │  │ Reporting                │  │
              │  │ Audit                    │  │
              │  │ Configuration            │  │
              │  └──────────────────────────┘  │
              └───────────────┬────────────────┘
                              │
                              ▼
                     ┌────────────────┐
                     │   PostgreSQL   │
                     └────────────────┘

                              ▲
                              │
                     ┌────────┴────────┐
                     │ Background      │
                     │ Worker          │
                     │                 │
                     │ Reminders       │
                     │ Outbox          │
                     │ Retry           │
                     │ Risk updates    │
                     └────────┬────────┘
                              │
                              ▼
                     ┌─────────────────┐
                     │ WhatsApp API    │
                     │ Provider        │
                     └─────────────────┘
```

This should be a **modular monolith**, not one huge coupled codebase.

---

## 2. Deployment Components

I would deploy two Go processes from the same codebase.

### API Server

Responsible for synchronous user-facing operations:

```text
Authentication
Create customer
Create purchase
Approve purchase
Create sale
Collect payment
Create expense
Read inventory
Read reports
Configure reminder rules
```

Example:

```text
business-api
```

### Worker

Responsible for asynchronous/background operations:

```text
WhatsApp notifications
Notification retries
Scheduled overdue scans
Customer risk classification
Outbox processing
Potential report aggregation
```

Example:

```text
business-worker
```

Both use the same PostgreSQL database.

This gives separation of workload without needing microservices.

---

## 3. Domain Architecture

The Go application should be divided into business modules.

```text
┌────────────────────────────────────────────┐
│                 API Layer                  │
│ HTTP handlers / middleware / validation    │
└───────────────────────┬────────────────────┘
                        │
┌───────────────────────▼────────────────────┐
│             Application Layer              │
│ Use cases / workflows / transactions       │
└───────────────────────┬────────────────────┘
                        │
┌───────────────────────▼────────────────────┐
│                Domain Layer                │
│ Business rules / entities / policies       │
└───────────────────────┬────────────────────┘
                        │
┌───────────────────────▼────────────────────┐
│            Infrastructure Layer            │
│ PostgreSQL / WhatsApp / files / logging    │
└────────────────────────────────────────────┘
```

Recommended modules:

| Module | Responsibility |
|---|---|
| IAM | Authentication, users, roles, permissions |
| Store | Store management |
| Customer | Customer identity and profile |
| Product | Company-wide product catalog |
| Procurement | Stock purchase workflow |
| Inventory | Stock quantities and movement ledger |
| Sales | Invoice creation and sale lifecycle |
| Credit | Due dates, outstanding amounts, aging |
| Payment | Collections and allocations |
| Expense | Store/company expenses |
| Notification | WhatsApp reminder management |
| Configuration | Business rules and thresholds |
| Reporting | Dashboard and finance reporting |
| Audit | Immutable audit trail |

---

## 4. Module Dependency Direction

Avoid allowing every module to access every other module arbitrarily.

A reasonable dependency graph is:

```text
                   IAM
                    │
                    ▼
                  Store

Product ──────► Procurement
                    │
                    ▼
                Inventory
                    ▲
                    │
Customer ───────► Sales
                    │
            ┌───────┴────────┐
            ▼                ▼
          Credit          Payment
            │                │
            └───────┬────────┘
                    ▼
              Notification

Expense ──────────────────────┐
Sales ────────────────────────┤
Payment ──────────────────────┤
Inventory ────────────────────┤
                              ▼
                          Reporting

All modules ───────────────► Audit
```

For example, the Notification module should not determine how much a customer owes.

Instead:

```text
Credit Module
   ↓
determines overdue invoice
   ↓
Notification Module
   ↓
sends reminder
```

---

## 5. PostgreSQL Architecture

Use **one PostgreSQL database** initially.

I would not create separate databases per module.

```text
business_management
```

You can organize tables logically using naming conventions or PostgreSQL schemas.

Example:

```text
iam.users
iam.roles
iam.permissions

inventory.inventory
inventory.inventory_movements

sales.invoices
sales.invoice_items

payment.payments
payment.payment_allocations
```

For a small engineering team, using the standard `public` schema with clean table names is also perfectly acceptable.

---

## 6. Core Data Model

At a high level:

```text
User
  │
  ├──── UserStore ───── Store
  │
  └──── Role / Permission


Customer
  │
  ├────────────── Invoice
  │                   │
  │                   ├── InvoiceItem ── Product
  │                   │
  │                   ├── PaymentAllocation
  │                   │        │
  │                   │        ▼
  │                   │      Payment
  │                   │
  │                   └── Notification
  │
  └── Customer Risk


Store
 │
 ├── Inventory ───── Product
 │        │
 │        └── InventoryMovement
 │
 ├── Purchase
 │      │
 │      ├── PurchaseItem
 │      ├── AdditionalCost
 │      └── Approval
 │
 ├── Invoice
 │
 └── Expense
```

---

## 7. Important Tables

The exact columns belong to LLD, but these are the important HLD-level tables.

### Identity

```text
users
roles
permissions
role_permissions
user_roles
user_stores
```

### Store

```text
stores
```

### Customer

```text
customers
customer_store_relationships
```

Customer phone number should have a unique constraint.

```sql
UNIQUE(phone_number)
```

If international phone numbers may be supported later, store normalized E.164 format.

Example:

```text
+919876543210
```

---

## 8. Product and Inventory

```text
products

inventory
inventory_movements
```

`inventory` represents the current snapshot:

```text
store_id
product_id
available_quantity
average_unit_cost
updated_at
```

Unique:

```text
(store_id, product_id)
```

`inventory_movements` represents history.

```text
movement_id
store_id
product_id
movement_type
quantity
unit_cost
reference_type
reference_id
created_at
created_by
```

Example:

```text
movement_type = PURCHASE
reference_id = PUR-1001
quantity = +100
```

or:

```text
movement_type = SALE
reference_id = INV-2010
quantity = -3
```

The movement ledger should be append-only.

---

## 9. Procurement Model

Recommended tables:

```text
purchases
purchase_items
purchase_additional_costs
purchase_approvals
```

Example:

```text
Purchase
────────────────────────
purchase_id
store_id
supplier_id
purchase_date
status
subtotal
gst_amount
additional_cost
total
created_by
approved_by
```

Items:

```text
PurchaseItem
────────────────────────
purchase_item_id
purchase_id
product_id
quantity
base_unit_cost
gst_amount
allocated_extra_cost
landed_unit_cost
```

---

## 10. Stock Approval Flow

The critical rule is:

> Creating a purchase request must not immediately modify inventory.

Flow:

```text
Store Admin
    │
    │ Create Purchase
    ▼
DRAFT
    │
    │ Submit
    ▼
PENDING_APPROVAL
    │
    │ Super Admin Approves
    ▼
APPROVED
    │
    ├── Calculate landed costs
    ├── Update weighted average cost
    ├── Increase inventory
    ├── Insert inventory movements
    └── Create audit event
```

All actions after approval should happen in **one PostgreSQL transaction**.

---

## 11. Purchase Approval Transaction

Conceptually:

```text
BEGIN TRANSACTION

1. Lock purchase
2. Verify status = PENDING_APPROVAL
3. Calculate additional cost allocation
4. Lock affected inventory rows
5. Update inventory quantities
6. Update weighted average cost
7. Insert InventoryMovement records
8. Change purchase status → APPROVED
9. Record approval
10. Write audit event

COMMIT
```

If any operation fails:

```text
ROLLBACK
```

So you never end up with:

```text
Purchase = APPROVED
Inventory = unchanged
```

or vice versa.

---

## 12. Weighted Average Inventory Cost

Suppose current stock:

```text
Quantity = 100
Average cost = ₹80
```

New purchase:

```text
Quantity = 50
Landed cost = ₹100
```

New weighted average:

```text
(100 × 80 + 50 × 100)
──────────────────────
        150

= ₹86.67
```

Store:

```text
available_quantity = 150
average_unit_cost = 86.67
```

When selling a product, capture the cost used at the time of sale.

Example:

```text
invoice_item.unit_cost = ₹86.67
```

Do **not** calculate historical profit using the latest inventory cost.

Otherwise a purchase next month would change last month's profit.

---

## 13. Sale Model

Recommended tables:

```text
invoices
invoice_items
```

Invoice:

```text
invoice_id
invoice_number
store_id
customer_id
subtotal
discount
gst
total_amount
paid_amount
outstanding_amount
payment_status
credit_status
due_date
status
created_by
created_at
```

Invoice items:

```text
invoice_item_id
invoice_id
product_id
quantity
sale_unit_price
unit_cost
gst
discount
line_total
```

Capturing both:

```text
sale_unit_price
unit_cost
```

is important.

Profit can then be derived from immutable historical information.

---

## 14. Sale Transaction

Sale creation is another critical atomic operation.

```text
BEGIN

1. Validate customer
2. Validate store
3. Lock inventory rows
4. Check quantity availability
5. Calculate invoice totals
6. Create invoice
7. Create invoice items
8. Reduce inventory
9. Create inventory movements
10. Create initial payment if amount_paid > 0
11. Calculate outstanding amount
12. Create audit event

COMMIT
```

Conceptually:

```text
Request
   │
   ▼
Validate products
   │
   ▼
Lock inventory
   │
   ▼
Stock sufficient?
   │
 ┌─┴───┐
 │     │
No    Yes
 │     │
Error  ▼
     Create Invoice
          │
          ▼
    Deduct Inventory
          │
          ▼
     Record Payment
          │
          ▼
     Create Credit
```

---

## 15. Preventing Overselling

Consider:

```text
Inventory = 1 unit
```

Two Store Admins simultaneously try to sell that unit.

Without concurrency protection:

```text
Admin A reads 1
Admin B reads 1

Admin A sells 1
Admin B sells 1
```

Now inventory effectively becomes `-1`.

Use PostgreSQL row locking.

Conceptually:

```sql
SELECT *
FROM inventory
WHERE store_id = $1
  AND product_id = $2
FOR UPDATE;
```

Then validate available quantity.

Alternatively:

```sql
UPDATE inventory
SET available_quantity = available_quantity - $qty
WHERE store_id = $store
AND product_id = $product
AND available_quantity >= $qty;
```

Then validate affected rows.

For multi-product sales, lock product rows in consistent order:

```text
Product ID ascending
```

to reduce deadlock risk.

---

## 16. Credit Architecture

I would **not create a separate mutable CustomerCredit table as the source of truth**.

The source of truth should be invoices and payment allocations.

Example:

```text
Invoice 1 outstanding = ₹5,000
Invoice 2 outstanding = ₹3,000
Invoice 3 outstanding = ₹0
```

Customer outstanding:

```text
SUM(outstanding_amount)

= ₹8,000
```

You may later maintain a cached summary for reporting, but the invoice data remains authoritative.

---

## 17. Payment Architecture

Tables:

```text
payments
payment_allocations
```

A payment represents money received.

Example:

```text
Payment

payment_id
customer_id
store_id
amount
payment_method
transaction_reference
collected_by
collected_at
status
```

Allocation tells us which invoice received the money.

```text
PaymentAllocation

payment_id
invoice_id
allocated_amount
```

This allows:

```text
One Payment
   │
   ├── ₹5,000 → Invoice 1
   └── ₹1,000 → Invoice 2
```

---

## 18. Payment Collection Flow

Example:

```text
Customer owes:

INV-1 → ₹5,000
INV-2 → ₹7,000

Customer pays ₹6,000
```

Default behavior:

```text
Oldest Invoice First

INV-1
₹5,000 applied

INV-2
₹1,000 applied
```

Transaction:

```text
BEGIN

1. Validate payment
2. Lock customer outstanding invoices
3. Create Payment
4. Allocate amount to invoices
5. Update invoice paid/outstanding values
6. Update payment statuses
7. Create audit record

COMMIT
```

---

## 19. Idempotency

This system handles money, therefore API retries must not duplicate transactions.

Consider:

```text
POST /payments

₹5,000 payment
```

Network timeout occurs after backend commits.

Frontend retries.

Without idempotency:

```text
Payment 1 → ₹5,000
Payment 2 → ₹5,000

Total = ₹10,000
```

which is incorrect.

Critical write APIs should accept:

```text
Idempotency-Key
```

Examples:

```text
POST /sales
POST /payments
POST /purchase/{id}/approve
```

You can maintain:

```text
idempotency_keys
```

containing:

```text
key
user_id
endpoint
request_hash
response
status
created_at
```

or place unique client request IDs directly on financial records.

---

## 20. Customer Risk Classification

Risk should be derived from the oldest outstanding overdue invoice.

Example:

```text
Customer

INV-A → 2 days overdue
INV-B → 10 days overdue
```

Maximum overdue:

```text
10 days
```

Configuration:

```text
0 days        → GREEN
1–6 days      → YELLOW
7+ days       → RED
```

Therefore:

```text
Customer = RED
```

The configuration should live in a system configuration table.

Example:

```text
customer_risk_config
```

---

## 21. Should Customer Risk Be Stored or Calculated?

Initially I recommend calculating it dynamically when retrieving a customer or report.

For example:

```text
MAX(current_date - due_date)
```

among outstanding invoices.

If reporting volume becomes large later, maintain a summary:

```text
customer_credit_summary
```

with:

```text
total_due
oldest_due_date
max_overdue_days
risk_status
```

updated asynchronously.

For V1, avoid unnecessary duplication.

---

## 22. WhatsApp Notification Architecture

Never send WhatsApp notifications directly inside a user transaction.

Bad:

```text
Create Payment
      ↓
Call WhatsApp API
      ↓
WhatsApp API takes 10 seconds
      ↓
User request blocked
```

Instead use asynchronous processing.

Recommended:

```text
Business Transaction
        │
        ▼
Store Notification Request
        │
        ▼
PostgreSQL
        │
        ▼
Background Worker
        │
        ▼
WhatsApp Provider
```

---

## 23. Transactional Outbox Pattern

I strongly recommend using an **Outbox Pattern**.

Suppose payment collection requires some future notification/event.

Inside the transaction:

```text
BEGIN

Create payment
Update invoice

INSERT INTO outbox_events (
    event_type = 'PAYMENT_RECEIVED',
    payload = ...
)

COMMIT
```

Then worker reads:

```text
outbox_events
```

and processes them asynchronously.

This prevents:

```text
Database committed
but application crashed before publishing event
```

You do **not need Kafka initially**.

PostgreSQL itself can act as the durable queue.

---

## 24. Outbox Table

Conceptually:

```text
outbox_events

id
event_type
aggregate_type
aggregate_id
payload
status
retry_count
available_at
created_at
processed_at
```

Statuses:

```text
PENDING
PROCESSING
COMPLETED
FAILED
```

Workers can fetch using:

```sql
FOR UPDATE SKIP LOCKED
```

This allows multiple workers to run simultaneously safely.

---

## 25. Automatic Reminder Flow

Suppose reminder configuration is:

```text
7 days
9 days
11 days
```

Worker periodically runs:

```text
Reminder Scheduler
        │
        ▼
Find overdue invoices
        │
        ▼
Determine overdue duration
        │
        ▼
Determine applicable reminder
        │
        ▼
Check if reminder already exists
        │
        ▼
Create Notification
        │
        ▼
WhatsApp Worker
        │
        ▼
Provider API
```

---

## 26. Prevent Duplicate Reminders

Create a uniqueness rule like:

```text
invoice_id
+
reminder_rule_id
```

Unique constraint:

```text
UNIQUE(invoice_id, reminder_rule_id)
```

Therefore even if scheduler runs multiple times:

```text
Invoice 100
Day 7 reminder
```

can only be created once.

This is much safer than relying only on application code.

---

## 27. Notification States

Recommended:

```text
PENDING
PROCESSING
SENT
DELIVERED
FAILED
```

Store:

```text
notification_id
customer_id
invoice_id
phone_number
type
trigger_type
reminder_rule_id
provider_message_id
status
retry_count
last_error
created_at
sent_at
```

---

## 28. WhatsApp Retry Strategy

Temporary provider errors should retry.

Example:

```text
Attempt 1 → immediately
Attempt 2 → +1 minute
Attempt 3 → +5 minutes
Attempt 4 → +30 minutes
```

After maximum retry:

```text
FAILED
```

Admin can see failed messages and optionally retry manually.

Use exponential backoff.

---

## 29. Scheduled Jobs

The worker should handle scheduled processes.

Example:

```text
business-worker

├── Reminder Scanner
├── Outbox Processor
├── WhatsApp Retry Worker
├── Risk Recalculation (if needed)
└── Reporting Aggregator (later)
```

Do not run critical cron logic independently inside every API replica.

If you deploy three API replicas:

```text
API 1
API 2
API 3
```

you don't want all three sending the same reminders.

Separate worker deployment solves most of this.

Where multiple workers exist, use PostgreSQL locking.

---

## 30. Expenses

Keep expenses simple:

```text
expenses

id
store_id NULLABLE
category_id
amount
expense_date
payment_method
vendor
description
receipt_url
created_by
status
created_at
```

Rule:

```text
store_id != NULL
→ Store Expense

store_id == NULL
→ Company Expense
```

If receipts are needed, store the file in object storage rather than PostgreSQL.

Postgres stores only:

```text
receipt_url
```

or object key.

---

## 31. Authentication Architecture

Recommended:

```text
Email
Password
    │
    ▼
Authentication API
    │
    ▼
Password hash verification
    │
    ▼
Access Token
```

Passwords:

```text
bcrypt / Argon2id
```

Never plaintext.

For a typical web admin application:

```text
Short-lived access token
+
Refresh token
```

is reasonable.

Alternatively secure server-side sessions are also valid.

---

## 32. Authorization Architecture

Use RBAC + store scope.

Example:

```text
User:
Amit

Role:
STORE_ADMIN

Stores:
STORE_101
```

Request:

```text
GET /stores/STORE_202/sales
```

Authentication succeeds, but authorization fails:

```text
403 Forbidden
```

because the user is not assigned to that store.

Authorization should be checked inside the backend.

Never trust frontend hiding buttons as authorization.

---

## 33. Permission Model

Possible permissions:

```text
CUSTOMER_CREATE
CUSTOMER_VIEW

SALE_CREATE
SALE_VIEW
SALE_CANCEL

PAYMENT_CREATE
PAYMENT_VIEW
PAYMENT_REVERSE

EXPENSE_CREATE
EXPENSE_VIEW

PURCHASE_CREATE
PURCHASE_APPROVE

INVENTORY_VIEW
INVENTORY_ADJUST

REMINDER_SEND

REPORT_VIEW_STORE
REPORT_VIEW_GLOBAL

USER_CREATE
USER_DISABLE

CONFIG_UPDATE
```

Then roles map to permissions.

This allows future roles without changing business code.

---

## 34. Audit Architecture

Audit log should be append-only.

```text
audit_logs
```

Fields:

```text
id
actor_user_id
store_id
action
entity_type
entity_id
old_value JSONB
new_value JSONB
ip_address
request_id
created_at
```

Example:

```text
action:
PURCHASE_APPROVED

entity:
PURCHASE

entity_id:
PUR-1201
```

Do not allow normal users to update/delete audit records.

---

## 35. Financial Transaction Reversal

Avoid:

```text
DELETE FROM payments
```

Instead:

```text
Original Payment
₹5,000

Reversal
-₹5,000
```

The system keeps both.

Same philosophy:

```text
Invoice
       ↓
Cancellation / Credit Note

Payment
       ↓
Payment Reversal

Inventory movement
       ↓
Reverse movement
```

This keeps accounting traceability.

---

## 36. Reporting Architecture

For V1, PostgreSQL can calculate reporting directly.

Examples:

```text
Today's Sales
Monthly Revenue
Outstanding Credit
Expenses
Gross Profit
Net Profit
Inventory Value
```

Use SQL aggregation.

Example:

```text
invoices
invoice_items
payments
expenses
```

should be the source of reports.

You do **not need a separate analytics database initially**.

---

## 37. Report Performance Strategy

Start with:

```text
PostgreSQL
+
proper indexes
+
SQL aggregation
```

When reports become slower:

```text
Materialized Views
```

can be introduced.

Example:

```text
monthly_store_financial_summary
```

Eventually:

```text
store_id
month
gross_sales
net_sales
cash_collected
cogs
gross_profit
expense
```

But don't build this in V1 unless necessary.

---

## 38. Important PostgreSQL Indexes

Examples:

```text
customers(phone_number)

inventory(store_id, product_id)

invoices(store_id, created_at)

invoices(customer_id, due_date)

invoices(customer_id, payment_status)

payments(customer_id, collected_at)

payments(store_id, collected_at)

inventory_movements(store_id, product_id, created_at)

notifications(status, available_at)

notifications(invoice_id, reminder_rule_id)

expenses(store_id, expense_date)
```

For reporting, date-based indexes will become important.

---

## 39. Database Constraints

Do not rely only on Go validation.

Important rules should also exist in PostgreSQL.

Examples:

```text
UNIQUE customer phone number

UNIQUE SKU

UNIQUE (store_id, product_id) inventory

quantity >= 0

payment amount > 0

expense amount > 0

outstanding amount >= 0
```

Foreign keys should also be used.

---

## 40. API Design

Use REST initially.

Base:

```text
/api/v1
```

Examples:

### Authentication

```http
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
```

### Customers

```http
POST /api/v1/customers
GET  /api/v1/customers/{id}
GET  /api/v1/customers?phone=...
GET  /api/v1/customers/{id}/invoices
GET  /api/v1/customers/{id}/payments
```

### Products

```http
POST /api/v1/products
GET  /api/v1/products
GET  /api/v1/products/{id}
```

### Procurement

```http
POST /api/v1/purchases
PUT  /api/v1/purchases/{id}
POST /api/v1/purchases/{id}/submit
POST /api/v1/purchases/{id}/approve
POST /api/v1/purchases/{id}/reject
```

### Inventory

```http
GET  /api/v1/stores/{storeId}/inventory
GET  /api/v1/stores/{storeId}/inventory/{productId}/movements
```

### Sales

```http
POST /api/v1/sales
GET  /api/v1/sales/{id}
GET  /api/v1/stores/{storeId}/sales
POST /api/v1/sales/{id}/cancel
```

### Payments

```http
POST /api/v1/payments
GET  /api/v1/payments/{id}
POST /api/v1/payments/{id}/reverse
```

### Expenses

```http
POST /api/v1/expenses
GET  /api/v1/expenses
```

### Reminders

```http
POST /api/v1/invoices/{id}/reminders
GET  /api/v1/customers/{id}/notifications
```

### Reporting

```http
GET /api/v1/reports/dashboard
GET /api/v1/reports/sales
GET /api/v1/reports/receivables
GET /api/v1/reports/profit
```

---

## 41. Request Tracing

Every HTTP request should have:

```text
request_id
```

Example:

```text
X-Request-ID
```

Logs:

```text
request_id
user_id
store_id
endpoint
duration
status
```

This helps trace:

```text
UI error
↓
API logs
↓
DB transaction
↓
Audit record
↓
Notification
```

---

## 42. Logging

Use structured logs in Go.

For example:

```json
{
  "level": "info",
  "request_id": "abc123",
  "user_id": "usr_100",
  "store_id": "store_10",
  "action": "create_sale",
  "invoice_id": "INV-1002"
}
```

Avoid logging sensitive information unnecessarily.

Particularly:

```text
Passwords
Tokens
Complete payment credentials
```

must never appear in logs.

---

## 43. Error Handling

API responses should have a consistent error contract.

Example:

```json
{
  "code": "INSUFFICIENT_STOCK",
  "message": "Requested quantity exceeds available stock",
  "details": {
    "product_id": "P100",
    "requested": 10,
    "available": 4
  }
}
```

Useful business error codes:

```text
CUSTOMER_ALREADY_EXISTS

INSUFFICIENT_STOCK

PURCHASE_ALREADY_APPROVED

INVALID_PURCHASE_STATE

PAYMENT_EXCEEDS_OUTSTANDING

INVOICE_ALREADY_PAID

UNAUTHORIZED_STORE_ACCESS

DUPLICATE_REQUEST
```

---

## 44. Configuration

Avoid hardcoding business rules.

Create configuration tables for:

```text
Credit Period

Reminder Rules

Risk Thresholds

Expense Categories

Product Categories

Payment Methods
```

Example:

```text
Reminder Rules

Rule 1 → 7 days
Rule 2 → 9 days
Rule 3 → 11 days
```

This means changing reminder timing doesn't require deployment.

---

## 45. File Storage

If you support:

```text
Expense receipts
Purchase bills
Invoice attachments
```

do not store large files as PostgreSQL blobs initially.

Architecture:

```text
Frontend
   │
   ▼
Backend
   │
   ▼
Object Storage
```

Examples could be S3-compatible storage.

PostgreSQL stores:

```text
object_key
file_name
content_type
```

---

## 46. Security

At minimum:

```text
HTTPS

Password hashing

Role-based authorization

Store-level authorization

Input validation

SQL parameterization

Secure token storage

Rate limiting on authentication

Audit logging

Secrets via environment/secrets manager

Database backups
```

Do not concatenate SQL strings from user inputs.

Using:

```text
pgx
```

or a suitable query layer with prepared/parameterized SQL is appropriate.

---

## 47. Deployment Architecture — Initial

For the initial production version:

```text
                    Internet
                       │
                       ▼
                Nginx / LB
                       │
           ┌───────────┴───────────┐
           │                       │
           ▼                       ▼
       Go API                  Go Worker
           │                       │
           └───────────┬───────────┘
                       │
                       ▼
                  PostgreSQL
                       │
                       ▼
                    Backup
```

WhatsApp worker communicates externally:

```text
Go Worker
   │
   ▼
WhatsApp Provider
```

---

## 48. Scaling

The API should be stateless.

Therefore when traffic increases:

```text
           Load Balancer
          /      |       \
         /       |        \
     API-1     API-2      API-3
        \        |        /
         \       |       /
            PostgreSQL
```

No user session should depend on a specific API instance.

---

## 49. PostgreSQL Connection Pooling

Each API replica should have a bounded connection pool.

For example:

```text
API-1 → 15 connections
API-2 → 15
API-3 → 15
Worker → 10
```

rather than every request opening a new DB connection.

Go's `pgxpool` is a good fit.

At larger scale, PgBouncer can be introduced.

---

## 50. PostgreSQL Scaling Path

Start:

```text
Single PostgreSQL Primary
```

Then, if needed:

```text
Primary
   │
   ├── Read Replica
   └── Read Replica
```

Writes:

```text
Sales
Payments
Inventory
```

go to primary.

Heavy reports can later go to read replicas.

Do not introduce this before necessary.

---

## 51. Backup and Recovery

Because financial data is involved, backup is mandatory.

At minimum:

```text
Daily database backup
+
PostgreSQL WAL / point-in-time recovery if available
```

Periodically verify that backups can actually be restored.

Backup without restore testing is not sufficient.

---

## 52. Suggested Go Project Structure

At a high level:

```text
cmd/
  api/
    main.go

  worker/
    main.go

internal/

  auth/
  users/
  stores/
  customers/
  products/

  procurement/
  inventory/
  sales/
  credit/
  payments/
  expenses/

  notifications/
  reporting/
  configuration/
  audit/

  platform/
    database/
    http/
    logging/
    middleware/
    whatsapp/
```

Inside each domain:

```text
sales/

  handler.go
  service.go
  repository.go
  model.go
```

or a cleaner layered variation.

Important principle:

```text
Handler
   ↓
Service / Use Case
   ↓
Repository
   ↓
PostgreSQL
```

Business rules shouldn't live in HTTP handlers.

---

## 53. Key Architecture Principle: Transactions Belong in Use Cases

For example:

```text
CreateSale()
```

should own the complete business transaction.

Not:

```text
SaleRepository creates invoice

InventoryService separately deducts stock

PaymentService separately inserts payment
```

with unrelated commits.

Instead:

```go
CreateSale(ctx, request)
    ↓
BEGIN
    ↓
Create invoice
Deduct inventory
Create movements
Create payment
Update outstanding
Create audit
    ↓
COMMIT
```

The service/use-case layer coordinates the transaction.

---

## 54. Major Domain Events

Even though we're not using Kafka, designing domain events is useful.

Examples:

```text
PurchaseSubmitted

PurchaseApproved

InventoryUpdated

SaleCreated

SaleCancelled

PaymentCollected

PaymentReversed

InvoiceOverdue

ReminderScheduled

ReminderSent
```

Initially these can be stored via the PostgreSQL Outbox.

Later, if architecture grows into microservices:

```text
PostgreSQL Outbox
        ↓
Kafka
        ↓
Independent Services
```

The domain model does not need to be redesigned.

---

## 55. Why I Would Not Use Kafka Initially

The system doesn't yet need:

```text
millions of events/sec
large numbers of independent consumers
cross-team microservices
event-stream analytics
```

Adding Kafka creates:

```text
Broker deployment
Schema management
Consumer lag
Retry strategy
DLQ
Monitoring
Operational maintenance
```

PostgreSQL + Outbox + workers is enough for this stage.

---

## 56. Why I Would Not Use Redis Initially

Redis could eventually help with:

```text
Caching
Distributed locks
Rate limiting
Session storage
```

but none is mandatory for V1.

PostgreSQL can initially handle:

```text
Configuration
Sessions
Outbox
Distributed job locking
```

Adding Redis only when there is a measurable need keeps the system simpler.

---

## 57. Failure Scenarios

### API crashes while creating sale

Because everything happens in one PostgreSQL transaction:

```text
Transaction incomplete
→ PostgreSQL rollback
```

No partial invoice.

### Payment request times out

Frontend retries with same:

```text
Idempotency-Key
```

Backend returns previous result.

No duplicate payment.

### WhatsApp provider is down

Notification remains:

```text
FAILED / RETRY_PENDING
```

Worker retries later.

Financial transaction is unaffected.

### Worker crashes

Outbox record remains unprocessed.

Another worker resumes processing.

### Two workers pick same notification

Use:

```text
FOR UPDATE SKIP LOCKED
```

plus unique constraints.

Only one processes the record.

### Two sales compete for same inventory

Use inventory row locking / atomic conditional update.

Only one can consume the remaining stock.

---

## 58. Data Consistency Classification

These should be **strongly consistent**:

```text
Inventory deduction

Sale creation

Payment recording

Payment allocation

Purchase approval

Inventory cost calculation
```

These can be **eventually consistent**:

```text
WhatsApp delivery

Customer risk cache

Dashboard summaries

Analytics

Notifications
```

This distinction is very important.

---

## 59. Reporting Definitions

The backend should expose financial values with clear definitions.

```text
Gross Sales
    =
Total invoice value

Revenue / Net Sales
    =
Gross sales
- discounts
- returns

Cash Collected
    =
Actual payment received

Outstanding Receivable
    =
Invoice total
- allocated payments

COGS
    =
Σ invoice_item.quantity × invoice_item.unit_cost

Gross Profit
    =
Revenue - COGS

Net Profit
    =
Gross Profit - Expenses
```

Do not mix:

```text
Revenue
```

and:

```text
Cash Collection
```

in database or API naming.

---

## 60. Important Architecture Decisions

For this system, the recommended decisions are:

| Concern | Recommendation |
|---|---|
| Architecture | Modular monolith |
| Language | Go |
| Database | PostgreSQL |
| API | REST |
| Database access | pgx / pgxpool |
| Transactions | PostgreSQL transactions |
| Inventory concurrency | Row locking / conditional updates |
| Background work | Separate Go worker |
| Messaging | PostgreSQL Outbox |
| Kafka | Not required initially |
| Redis | Not required initially |
| Inventory costing | Weighted average |
| Credit source of truth | Invoice-level outstanding |
| Payment allocation | Oldest-first with explicit override |
| Customer uniqueness | WhatsApp/mobile number |
| Reporting | PostgreSQL queries initially |
| Audit | Append-only audit log |
| Notifications | Async worker |
| Files | Object storage |
| Horizontal scaling | Stateless API replicas |

---

## 61. Final HLD

```text
                        ┌───────────────────┐
                        │     Web Admin     │
                        └─────────┬─────────┘
                                  │
                               HTTPS
                                  │
                                  ▼
                        ┌───────────────────┐
                        │  Nginx / LB       │
                        └─────────┬─────────┘
                                  │
                  ┌───────────────┴───────────────┐
                  │                               │
                  ▼                               ▼
        ┌───────────────────┐          ┌─────────────────────┐
        │      Go API       │          │      Go Worker      │
        │                   │          │                     │
        │ Auth              │          │ Scheduler           │
        │ Store             │          │ Outbox Processor    │
        │ Customer          │          │ WhatsApp Sender     │
        │ Product           │          │ Retry Processor     │
        │ Procurement       │          │                     │
        │ Inventory         │          └──────────┬──────────┘
        │ Sales             │                     │
        │ Credit            │                     ▼
        │ Payment           │           ┌───────────────────┐
        │ Expense           │           │ WhatsApp Provider │
        │ Reporting         │           └───────────────────┘
        │ Configuration     │
        │ Audit             │
        └──────────┬────────┘
                   │
                   │
                   ▼
          ┌────────────────────┐
          │     PostgreSQL     │
          │                    │
          │ Operational Data   │
          │ Financial Data     │
          │ Inventory Ledger   │
          │ Outbox Queue       │
          │ Audit Logs         │
          │ Configuration      │
          └──────────┬─────────┘
                     │
                     ▼
                DB Backups
```

The most important architecture choice is **keeping financial and inventory operations inside PostgreSQL transactions while pushing non-critical asynchronous work such as WhatsApp to the worker/outbox system**. This gives strong consistency where money and stock are involved without introducing the complexity of microservices too early.

The natural next step from this HLD is the **database/ER design**, especially:

```text
Purchase → Inventory → Sale → COGS
```

and:

```text
Customer → Invoice → Payment → Allocation → Outstanding
```

Once those tables and relationships are fixed, the API contracts and Go interfaces become significantly easier to define.
