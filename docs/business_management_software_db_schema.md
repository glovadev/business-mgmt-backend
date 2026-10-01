# Database Schema — Business Management Software

**Backend:** Golang  
**Database:** PostgreSQL

This document defines the proposed V1 database schema for the Business Management Software, including tables, fields, constraints, indexes, source-of-truth decisions, transaction boundaries, and important database invariants.

---

## 1. Database Conventions

Use these conventions consistently:

```text
Primary keys       UUID
Money              NUMERIC(18,2)
Unit cost           NUMERIC(18,4)
Quantity            NUMERIC(18,3)
Percentage          NUMERIC(5,2)
Timestamp           TIMESTAMPTZ
Business date       DATE
Phone               VARCHAR(20), normalized E.164
Currency            CHAR(3), default INR
```

Examples:

```text
₹1,250.50        → NUMERIC(18,2)

Unit cost ₹12.3456
                 → NUMERIC(18,4)

12.500 KG        → NUMERIC(18,3)
```

Even if products currently use integer quantities, fractional quantities keep the model flexible.

Recommended conventions:

- Use UUID primary keys.
- Use `TIMESTAMPTZ` for timestamps.
- Use `NUMERIC` for monetary values.
- Use normal `VARCHAR/TEXT + CHECK` constraints instead of PostgreSQL enums for statuses.
- Normalize phone numbers before storage.
- Normalize emails to lowercase or use PostgreSQL `CITEXT`.

---

## 2. Schema Overview

```text
IAM
├── users
├── roles
├── permissions
├── user_roles
├── role_permissions
├── user_stores
└── user_sessions

Store
└── stores

Customer
├── customers
└── customer_store_relationships

Product
├── product_categories
├── products
└── suppliers

Procurement
├── purchases
├── purchase_items
├── purchase_additional_costs
└── purchase_approvals

Inventory
├── inventory
└── inventory_movements

Sales
├── invoices
└── invoice_items

Payments
├── payments
└── payment_allocations

Expense
├── expense_categories
└── expenses

Configuration
├── company_settings
└── reminder_rules

Notifications
├── notification_templates
└── notifications

Infrastructure
├── idempotency_keys
├── outbox_events
├── audit_logs
└── document_counters
```

---

## 3. `stores`

| Field | Type | Notes |
|---|---|---|
| id | UUID PK | Internal ID |
| code | VARCHAR(30) | Human-friendly store code |
| name | VARCHAR(150) | Store name |
| phone | VARCHAR(20) | Optional |
| gstin | VARCHAR(20) | Optional |
| address_line1 | VARCHAR(255) | |
| address_line2 | VARCHAR(255) | Nullable |
| city | VARCHAR(100) | |
| state | VARCHAR(100) | |
| postal_code | VARCHAR(20) | |
| country_code | CHAR(2) | `IN` |
| status | VARCHAR(20) | ACTIVE / INACTIVE |
| created_by | UUID FK users | Nullable for bootstrap |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

Constraints:

```sql
UNIQUE(code)

CHECK (status IN ('ACTIVE', 'INACTIVE'))
```

Indexes:

```sql
CREATE INDEX idx_stores_status
ON stores(status);
```

---

## 4. `users`

| Field | Type |
|---|---|
| id | UUID PK |
| email | VARCHAR(255) |
| password_hash | TEXT |
| full_name | VARCHAR(150) |
| phone | VARCHAR(20) nullable |
| status | VARCHAR(20) |
| last_login_at | TIMESTAMPTZ nullable |
| created_at | TIMESTAMPTZ |
| updated_at | TIMESTAMPTZ |

Constraints:

```sql
UNIQUE(LOWER(email))

CHECK (
    status IN ('ACTIVE', 'INACTIVE', 'LOCKED')
)
```

Ideally use PostgreSQL `CITEXT` for email or normalize email to lowercase before storage.

---

## 5. RBAC Tables

### `roles`

```text
id
code
name
description
created_at
```

Example codes:

```text
SUPER_ADMIN
FINANCE_ADMIN
STORE_ADMIN
```

Constraint:

```sql
UNIQUE(code)
```

### `permissions`

```text
id
code
description
created_at
```

Examples:

```text
SALE_CREATE
SALE_VIEW

PAYMENT_CREATE
PAYMENT_REVERSE

PURCHASE_CREATE
PURCHASE_APPROVE

INVENTORY_VIEW
INVENTORY_ADJUST

EXPENSE_CREATE

REMINDER_SEND

REPORT_VIEW_GLOBAL
```

Constraint:

```sql
UNIQUE(code)
```

### `user_roles`

```text
user_id UUID FK
role_id UUID FK
created_at TIMESTAMPTZ
```

Primary key:

```sql
PRIMARY KEY(user_id, role_id)
```

### `role_permissions`

```text
role_id
permission_id
```

Primary key:

```sql
PRIMARY KEY(role_id, permission_id)
```

---

## 6. `user_stores`

Controls store scope for users.

```text
user_id
store_id
created_at
```

Primary key:

```sql
PRIMARY KEY(user_id, store_id)
```

Index:

```sql
CREATE INDEX idx_user_stores_store
ON user_stores(store_id);
```

Example:

```text
Store Admin
→ Store A
```

Finance Admin and Super Admin may use global permission instead of requiring entries for every store.

---

## 7. `user_sessions`

If using refresh-token authentication:

```text
id UUID
user_id UUID
refresh_token_hash TEXT
device_info TEXT
ip_address INET
expires_at TIMESTAMPTZ
revoked_at TIMESTAMPTZ NULL
created_at TIMESTAMPTZ
```

Important:

Store only a **hash of the refresh token**, not the token itself.

Indexes:

```sql
CREATE INDEX idx_user_sessions_user
ON user_sessions(user_id);

CREATE INDEX idx_user_sessions_expiry
ON user_sessions(expires_at)
WHERE revoked_at IS NULL;
```

---

## 8. `customers`

| Field | Type |
|---|---|
| id | UUID PK |
| phone_number | VARCHAR(20) |
| name | VARCHAR(150) |
| alternate_phone | VARCHAR(20) nullable |
| email | VARCHAR(255) nullable |
| address | TEXT nullable |
| gstin | VARCHAR(20) nullable |
| whatsapp_opt_in | BOOLEAN |
| whatsapp_opted_out_at | TIMESTAMPTZ nullable |
| status | VARCHAR(20) |
| notes | TEXT nullable |
| created_by | UUID FK users |
| created_at | TIMESTAMPTZ |
| updated_at | TIMESTAMPTZ |

Phone should be normalized.

Example:

```text
9876543210
```

becomes:

```text
+919876543210
```

Constraints:

```sql
UNIQUE(phone_number)

CHECK(status IN ('ACTIVE', 'INACTIVE'))
```

Indexes:

```sql
CREATE UNIQUE INDEX uq_customers_phone
ON customers(phone_number);

CREATE INDEX idx_customers_name
ON customers(name);
```

If customer search by name becomes important, PostgreSQL trigram indexing can be introduced later.

---

## 9. `customer_store_relationships`

Useful because a Store Admin may add a customer before a sale occurs.

```text
customer_id
store_id
first_added_by
created_at
```

Primary key:

```sql
PRIMARY KEY(customer_id, store_id)
```

Index:

```sql
CREATE INDEX idx_customer_store_store
ON customer_store_relationships(store_id, customer_id);
```

---

## 10. `product_categories`

```text
id UUID
name VARCHAR(100)
description TEXT
status VARCHAR(20)
created_at TIMESTAMPTZ
```

Constraint:

```sql
UNIQUE(name)
```

---

## 11. `products`

Product definition is company-wide.

| Field | Type |
|---|---|
| id | UUID PK |
| sku | VARCHAR(100) |
| name | VARCHAR(200) |
| category_id | UUID FK |
| description | TEXT |
| hsn_code | VARCHAR(20) nullable |
| unit_of_measure | VARCHAR(30) |
| default_selling_price | NUMERIC(18,2) |
| gst_rate | NUMERIC(5,2) |
| status | VARCHAR(20) |
| created_by | UUID |
| created_at | TIMESTAMPTZ |
| updated_at | TIMESTAMPTZ |

Examples:

```text
unit_of_measure

PIECE
KG
LITRE
BOX
BAG
BUCKET
```

Constraints:

```sql
UNIQUE(sku)

CHECK(default_selling_price >= 0)

CHECK(gst_rate >= 0 AND gst_rate <= 100)

CHECK(status IN ('ACTIVE', 'INACTIVE'))
```

Indexes:

```sql
CREATE INDEX idx_products_category
ON products(category_id);

CREATE INDEX idx_products_status
ON products(status);
```

---

## 12. `suppliers`

Even though supplier management is not a major PRD feature, procurement should reference suppliers.

```text
id
name
phone
email
gstin
address
status
created_at
updated_at
```

Index:

```sql
CREATE INDEX idx_suppliers_name
ON suppliers(name);
```

---

## 13. `purchases`

Represents one procurement batch.

| Field | Type |
|---|---|
| id | UUID PK |
| purchase_number | VARCHAR(50) |
| store_id | UUID FK |
| supplier_id | UUID FK nullable |
| supplier_invoice_number | VARCHAR(100) nullable |
| purchase_date | DATE |
| status | VARCHAR(30) |
| subtotal | NUMERIC(18,2) |
| gst_amount | NUMERIC(18,2) |
| additional_cost_total | NUMERIC(18,2) |
| grand_total | NUMERIC(18,2) |
| notes | TEXT nullable |
| created_by | UUID |
| submitted_at | TIMESTAMPTZ nullable |
| approved_by | UUID nullable |
| approved_at | TIMESTAMPTZ nullable |
| created_at | TIMESTAMPTZ |
| updated_at | TIMESTAMPTZ |

Statuses:

```text
DRAFT
PENDING_APPROVAL
APPROVED
REJECTED
CANCELLED
```

Constraint:

```sql
UNIQUE(purchase_number)
```

Indexes:

```sql
CREATE INDEX idx_purchases_store_date
ON purchases(store_id, purchase_date DESC);

CREATE INDEX idx_purchases_pending
ON purchases(status, created_at)
WHERE status = 'PENDING_APPROVAL';
```

The partial index is useful because Super Admin will frequently query pending purchases.

---

## 14. `purchase_items`

| Field | Type |
|---|---|
| id | UUID PK |
| purchase_id | UUID FK |
| product_id | UUID FK |
| quantity | NUMERIC(18,3) |
| base_unit_cost | NUMERIC(18,4) |
| base_amount | NUMERIC(18,2) |
| gst_rate | NUMERIC(5,2) |
| gst_amount | NUMERIC(18,2) |
| allocated_extra_cost | NUMERIC(18,2) |
| landed_unit_cost | NUMERIC(18,4) |
| total_amount | NUMERIC(18,2) |
| created_at | TIMESTAMPTZ |

Constraints:

```sql
CHECK(quantity > 0)

CHECK(base_unit_cost >= 0)

CHECK(landed_unit_cost >= 0)
```

Recommended V1 constraint:

```sql
UNIQUE(purchase_id, product_id)
```

This simplifies procurement because the same product occurs once per purchase.

Index:

```sql
CREATE INDEX idx_purchase_items_product
ON purchase_items(product_id);
```

---

## 15. `purchase_additional_costs`

Instead of fixed columns for transport, loading, insurance, and similar costs, keep these dynamic.

```text
id
purchase_id
cost_type
description
amount
created_at
```

Example:

```text
TRANSPORTATION     ₹5,000
LOADING            ₹800
TRAVEL             ₹2,500
INSURANCE          ₹1,000
OTHER              ₹500
```

Constraint:

```sql
CHECK(amount >= 0)
```

Index:

```sql
CREATE INDEX idx_purchase_costs_purchase
ON purchase_additional_costs(purchase_id);
```

---

## 16. `purchase_approvals`

Provides explicit approval history.

```text
id
purchase_id
action
actor_user_id
remarks
created_at
```

Actions:

```text
SUBMITTED
APPROVED
REJECTED
CANCELLED
```

Index:

```sql
CREATE INDEX idx_purchase_approvals_purchase
ON purchase_approvals(purchase_id, created_at);
```

Even though `purchases.approved_by` exists for quick access, this table is the historical trail.

---

## 17. `inventory`

This is the **current inventory snapshot**.

| Field | Type |
|---|---|
| store_id | UUID FK |
| product_id | UUID FK |
| available_quantity | NUMERIC(18,3) |
| average_unit_cost | NUMERIC(18,4) |
| updated_at | TIMESTAMPTZ |

Primary key:

```sql
PRIMARY KEY(store_id, product_id)
```

Constraints:

```sql
CHECK(available_quantity >= 0)

CHECK(average_unit_cost >= 0)
```

This table should remain small and fast compared with the movement ledger.

---

## 18. `inventory_movements`

This is the **inventory source-of-history**.

Never delete these records.

| Field | Type |
|---|---|
| id | UUID PK |
| store_id | UUID |
| product_id | UUID |
| movement_type | VARCHAR(30) |
| quantity_change | NUMERIC(18,3) |
| unit_cost | NUMERIC(18,4) nullable |
| total_cost | NUMERIC(18,2) nullable |
| source_type | VARCHAR(30) |
| source_id | UUID nullable |
| source_line_id | UUID nullable |
| reason | TEXT nullable |
| created_by | UUID |
| created_at | TIMESTAMPTZ |

Movement types:

```text
PURCHASE
SALE
SALE_RETURN
PURCHASE_RETURN
ADJUSTMENT_IN
ADJUSTMENT_OUT
TRANSFER_IN
TRANSFER_OUT
REVERSAL
```

Example:

```text
Purchase approved

quantity_change = +100
```

Sale:

```text
quantity_change = -4
```

Indexes:

```sql
CREATE INDEX idx_inventory_movements_product_history
ON inventory_movements(
    store_id,
    product_id,
    created_at DESC
);

CREATE INDEX idx_inventory_movements_source
ON inventory_movements(
    source_type,
    source_id
);
```

For transaction-created movements:

```sql
CREATE UNIQUE INDEX uq_inventory_movement_source_line
ON inventory_movements(
    movement_type,
    source_line_id
)
WHERE source_line_id IS NOT NULL;
```

This helps prevent a purchase or invoice line from modifying inventory twice.

---

## 19. `invoices`

Represents the sale.

| Field | Type |
|---|---|
| id | UUID PK |
| invoice_number | VARCHAR(50) |
| store_id | UUID FK |
| customer_id | UUID FK |
| invoice_date | DATE |
| subtotal | NUMERIC(18,2) |
| discount_amount | NUMERIC(18,2) |
| gst_amount | NUMERIC(18,2) |
| total_amount | NUMERIC(18,2) |
| paid_amount | NUMERIC(18,2) |
| outstanding_amount | NUMERIC(18,2) |
| payment_status | VARCHAR(30) |
| due_date | DATE nullable |
| status | VARCHAR(30) |
| created_by | UUID |
| created_at | TIMESTAMPTZ |
| updated_at | TIMESTAMPTZ |

Payment status:

```text
UNPAID
PARTIALLY_PAID
PAID
```

Invoice status:

```text
CONFIRMED
CANCELLED
```

Avoid storing a separate `credit_status` because values such as:

```text
NOT_DUE
OVERDUE
SETTLED
```

can be calculated from:

```text
outstanding_amount
due_date
current_date
```

This prevents stale status values.

Constraints:

```sql
UNIQUE(invoice_number)

CHECK(total_amount >= 0)

CHECK(paid_amount >= 0)

CHECK(outstanding_amount >= 0)

CHECK(paid_amount <= total_amount)

CHECK (
    outstanding_amount = 0
    OR due_date IS NOT NULL
)
```

---

## 20. Important Invoice Indexes

Store sales:

```sql
CREATE INDEX idx_invoices_store_date
ON invoices(store_id, invoice_date DESC);
```

Customer history:

```sql
CREATE INDEX idx_invoices_customer_date
ON invoices(customer_id, invoice_date DESC);
```

Outstanding invoices:

```sql
CREATE INDEX idx_invoices_customer_outstanding
ON invoices(customer_id, due_date)
WHERE outstanding_amount > 0
  AND status = 'CONFIRMED';
```

Overdue reminder queries:

```sql
CREATE INDEX idx_invoices_due
ON invoices(due_date)
WHERE outstanding_amount > 0
  AND status = 'CONFIRMED';
```

The partial overdue index is especially important for the scheduled reminder worker.

---

## 21. `invoice_items`

| Field | Type |
|---|---|
| id | UUID PK |
| invoice_id | UUID |
| product_id | UUID |
| quantity | NUMERIC(18,3) |
| unit_price | NUMERIC(18,2) |
| discount_amount | NUMERIC(18,2) |
| taxable_amount | NUMERIC(18,2) |
| gst_rate | NUMERIC(5,2) |
| gst_amount | NUMERIC(18,2) |
| line_total | NUMERIC(18,2) |
| unit_cost | NUMERIC(18,4) |
| created_at | TIMESTAMPTZ |

The critical field is:

```text
unit_cost
```

Suppose the inventory average cost at sale time is:

```text
₹86.67
```

Store:

```text
invoice_items.unit_cost = 86.67
```

permanently.

Historical COGS is then:

```text
quantity × unit_cost
```

Future inventory purchases cannot alter historical profit.

Indexes:

```sql
CREATE INDEX idx_invoice_items_invoice
ON invoice_items(invoice_id);

CREATE INDEX idx_invoice_items_product
ON invoice_items(product_id);
```

---

## 22. `payments`

A payment is the **actual money received**.

| Field | Type |
|---|---|
| id | UUID PK |
| payment_number | VARCHAR(50) |
| customer_id | UUID |
| amount | NUMERIC(18,2) |
| payment_method | VARCHAR(30) |
| transaction_reference | VARCHAR(150) nullable |
| collection_store_id | UUID nullable |
| collected_by | UUID |
| collected_at | TIMESTAMPTZ |
| status | VARCHAR(30) |
| reversal_of_payment_id | UUID nullable |
| notes | TEXT |
| created_at | TIMESTAMPTZ |

Payment methods:

```text
CASH
UPI
BANK_TRANSFER
```

Status:

```text
POSTED
REVERSED
```

Constraint:

```sql
CHECK(amount > 0)
```

Unique:

```sql
UNIQUE(payment_number)
```

---

## 23. Why `collection_store_id` Is Nullable

Example:

```text
Customer Rahul

Store A Invoice = ₹5,000
Store B Invoice = ₹3,000
```

Finance Admin receives:

```text
₹8,000
```

centrally.

The payment itself may not belong to either physical store.

Therefore:

```text
payment.customer_id = Rahul

payment.collection_store_id = NULL
```

Then:

```text
PaymentAllocation

₹5,000 → Store A invoice
₹3,000 → Store B invoice
```

For accounting by **sale store**, derive from invoice allocations.

For reporting by **where cash was physically collected**, use `collection_store_id`.

---

## 24. Payment Indexes

```sql
CREATE INDEX idx_payments_customer_date
ON payments(customer_id, collected_at DESC);

CREATE INDEX idx_payments_collector_date
ON payments(collected_by, collected_at DESC);

CREATE INDEX idx_payments_collection_store
ON payments(collection_store_id, collected_at DESC)
WHERE collection_store_id IS NOT NULL;
```

---

## 25. `payment_allocations`

This table maps received money to invoices.

| Field | Type |
|---|---|
| id | UUID |
| payment_id | UUID |
| invoice_id | UUID |
| allocated_amount | NUMERIC(18,2) |
| created_at | TIMESTAMPTZ |

Constraint:

```sql
CHECK(allocated_amount > 0)
```

Unique:

```sql
UNIQUE(payment_id, invoice_id)
```

Indexes:

```sql
CREATE INDEX idx_payment_allocations_invoice
ON payment_allocations(invoice_id);

CREATE INDEX idx_payment_allocations_payment
ON payment_allocations(payment_id);
```

Example:

```text
Payment ₹6,000
   │
   ├── INV-1 ₹5,000
   └── INV-2 ₹1,000
```

---

## 26. Payment Invariants

These cannot all be represented through simple `CHECK` constraints.

The Go transaction must guarantee:

```text
SUM(payment_allocations)
<= payment.amount
```

and:

```text
allocated_amount
<= invoice.outstanding_amount
```

The application should lock the relevant invoices before allocation.

---

## 27. Source of Truth for Outstanding Amount

Conceptually, the accounting source of truth is:

```text
Invoice Total
-
Successful Payment Allocations
```

However, calculating that repeatedly is unnecessarily expensive.

Therefore keep:

```text
invoices.paid_amount
invoices.outstanding_amount
```

as **transactionally maintained denormalized fields**.

Meaning:

```text
Payment allocation inserted
+
Invoice paid_amount updated
+
Invoice outstanding_amount updated
```

must happen in the **same database transaction**.

This provides:

```text
correct accounting
+
fast queries
```

---

## 28. `expense_categories`

```text
id
name
scope
status
created_at
```

Scope:

```text
STORE
COMPANY
BOTH
```

Examples:

```text
Rent
Electricity
Tea / Coffee
Internet
Marketing
Legal
Travel
```

---

## 29. `expenses`

| Field | Type |
|---|---|
| id | UUID |
| expense_number | VARCHAR(50) |
| store_id | UUID nullable |
| category_id | UUID |
| expense_date | DATE |
| amount | NUMERIC(18,2) |
| gst_amount | NUMERIC(18,2) |
| payment_method | VARCHAR(30) nullable |
| vendor | VARCHAR(200) nullable |
| description | TEXT |
| receipt_object_key | TEXT nullable |
| status | VARCHAR(20) |
| created_by | UUID |
| reversal_of_expense_id | UUID nullable |
| created_at | TIMESTAMPTZ |

Rule:

```text
store_id != NULL
→ Store Expense

store_id = NULL
→ Company Expense
```

Indexes:

```sql
CREATE INDEX idx_expenses_store_date
ON expenses(store_id, expense_date DESC);

CREATE INDEX idx_expenses_category_date
ON expenses(category_id, expense_date DESC);
```

---

## 30. `company_settings`

For V1, one row can hold important company-level business policies.

Fields:

```text
id

company_name

currency
timezone

default_credit_days

yellow_after_overdue_days
red_after_overdue_days

additional_cost_allocation_method

include_gst_in_landed_cost

created_at
updated_at
updated_by
```

Example:

```text
currency = INR

timezone = Asia/Kolkata

default_credit_days = 15

yellow_after_overdue_days = 1

red_after_overdue_days = 7

additional_cost_allocation_method =
PURCHASE_VALUE

include_gst_in_landed_cost = false
```

Constraint:

```text
red_after_overdue_days >
yellow_after_overdue_days
```

This keeps business policy out of source code.

---

## 31. `reminder_rules`

Instead of fixed Day-7, Day-9, and Day-11 rules, store them in data.

```text
id
sequence_number
days_after_due
template_id
status
created_by
created_at
updated_at
```

Example:

| Sequence | Days |
|---:|---:|
| 1 | 7 |
| 2 | 9 |
| 3 | 11 |

Constraints:

```sql
UNIQUE(sequence_number)

UNIQUE(days_after_due)
```

---

## 32. `notification_templates`

Useful because WhatsApp Business commonly works with registered templates.

```text
id
code
provider_template_id
language_code
body_preview
status
created_at
updated_at
```

Example:

```text
PAYMENT_DUE_REMINDER
```

---

## 33. `notifications`

| Field | Type |
|---|---|
| id | UUID |
| customer_id | UUID |
| invoice_id | UUID nullable |
| reminder_rule_id | UUID nullable |
| phone_number | VARCHAR(20) |
| trigger_type | VARCHAR(20) |
| notification_type | VARCHAR(50) |
| template_id | UUID |
| payload | JSONB |
| provider | VARCHAR(50) |
| provider_message_id | VARCHAR(255) nullable |
| status | VARCHAR(20) |
| retry_count | INTEGER |
| next_retry_at | TIMESTAMPTZ nullable |
| last_error | TEXT nullable |
| triggered_by | UUID nullable |
| created_at | TIMESTAMPTZ |
| sent_at | TIMESTAMPTZ nullable |
| delivered_at | TIMESTAMPTZ nullable |

Trigger type:

```text
AUTOMATIC
MANUAL
```

Status:

```text
PENDING
PROCESSING
SENT
DELIVERED
FAILED
```

---

## 34. Prevent Duplicate Automatic Reminders

Important unique index:

```sql
CREATE UNIQUE INDEX uq_automatic_invoice_reminder
ON notifications(invoice_id, reminder_rule_id)
WHERE trigger_type = 'AUTOMATIC'
  AND reminder_rule_id IS NOT NULL;
```

Therefore:

```text
INV-1001
Day-7 rule
```

can only be automatically created once.

Manual reminders remain unaffected.

---

## 35. Notification Worker Index

The worker frequently queries pending or retryable notifications.

```sql
CREATE INDEX idx_notifications_worker
ON notifications(status, next_retry_at)
WHERE status IN ('PENDING', 'FAILED');
```

---

## 36. `idempotency_keys`

Use for critical write APIs such as:

```text
POST /sales
POST /payments
POST /purchases/{id}/approve
```

Fields:

```text
id
idempotency_key
user_id
http_method
endpoint
request_hash
status
response_code
response_body JSONB
created_at
expires_at
```

Constraint:

```sql
UNIQUE(user_id, idempotency_key)
```

The request hash prevents reusing the same key for different requests.

Example:

```text
Same Idempotency-Key

Request 1 = ₹5,000
Request 2 = ₹10,000
```

Request 2 should be rejected.

---

## 37. `outbox_events`

Used for asynchronous work.

Fields:

```text
id UUID

event_type VARCHAR(100)

aggregate_type VARCHAR(50)

aggregate_id UUID

payload JSONB

status VARCHAR(20)

retry_count INTEGER

available_at TIMESTAMPTZ

locked_at TIMESTAMPTZ nullable

locked_by VARCHAR(100) nullable

last_error TEXT nullable

created_at TIMESTAMPTZ

processed_at TIMESTAMPTZ nullable
```

Statuses:

```text
PENDING
PROCESSING
COMPLETED
FAILED
```

Worker index:

```sql
CREATE INDEX idx_outbox_pending
ON outbox_events(available_at, created_at)
WHERE status IN ('PENDING', 'FAILED');
```

Worker fetch pattern:

```sql
SELECT ...
FROM outbox_events
WHERE status = 'PENDING'
  AND available_at <= NOW()
ORDER BY created_at
FOR UPDATE SKIP LOCKED
LIMIT 100;
```

---

## 38. `audit_logs`

This table should be append-only.

Fields:

```text
id UUID

actor_user_id UUID

store_id UUID nullable

action VARCHAR(100)

entity_type VARCHAR(100)

entity_id UUID

old_value JSONB nullable

new_value JSONB nullable

request_id VARCHAR(100)

ip_address INET nullable

created_at TIMESTAMPTZ
```

Example:

```text
actor_user_id:
USR-100

action:
PAYMENT_CREATED

entity_type:
PAYMENT

entity_id:
PAY-1002
```

Indexes:

```sql
CREATE INDEX idx_audit_entity
ON audit_logs(entity_type, entity_id, created_at DESC);

CREATE INDEX idx_audit_actor
ON audit_logs(actor_user_id, created_at DESC);

CREATE INDEX idx_audit_store
ON audit_logs(store_id, created_at DESC);
```

Do not immediately add GIN indexes on:

```text
old_value
new_value
```

unless JSON search is truly needed.

---

## 39. `document_counters`

Used for human-readable document numbers.

Examples:

```text
INV-BRP-2026-000123

PUR-BRP-2026-000021

PAY-2026-000842
```

Do not generate these using:

```text
SELECT COUNT(*) + 1
```

because concurrent requests may generate duplicates.

Fields:

```text
store_id nullable
document_type
financial_year
current_value
updated_at
```

Primary key:

```sql
PRIMARY KEY(
    store_id,
    document_type,
    financial_year
)
```

Generate the next value inside a transaction while locking the counter row.

---

## 40. Entity Relationships

Core financial relationship:

```text
Customer
   │
   ├──────────── Invoice
   │                │
   │                ├── InvoiceItem
   │                │       │
   │                │       ▼
   │                │     Product
   │                │
   │                ▼
   │        PaymentAllocation
   │                ▲
   │                │
   └──────────── Payment
```

Inventory relationship:

```text
Supplier
    │
    ▼
Purchase
    │
    ├── PurchaseItem
    │        │
    │        ▼
    │      Product
    │
    ▼
InventoryMovement
    │
    ▼
Inventory Snapshot
```

Sale relationship:

```text
Invoice
   │
   ▼
InvoiceItem
   │
   ▼
InventoryMovement (-)
```

Complete trace:

```text
Purchase
   ↓
Inventory
   ↓
Sale
   ↓
COGS
   ↓
Profit
```

---

## 41. Foreign-Key Deletion Strategy

For financial systems, avoid cascading deletes for transaction data.

Do not use:

```sql
ON DELETE CASCADE
```

for chains such as:

```text
customer
→ invoices
→ payments
```

A customer deletion should never cause financial records to disappear.

Prefer:

```text
Master data
→ ACTIVE / INACTIVE

Financial transactions
→ never delete
```

Examples:

```text
Product inactive

Customer inactive

Store inactive

User inactive
```

Historical records remain.

---

## 42. Transaction Boundaries

### Create Sale

One transaction:

```text
BEGIN

Lock inventory rows

Create invoice

Create invoice_items

Update inventory

Create inventory_movements

Create initial payment

Create payment allocation

Update invoice paid/outstanding

Create audit log

COMMIT
```

### Approve Purchase

```text
BEGIN

Lock purchase

Validate state

Lock inventory

Calculate landed cost

Update inventory

Create inventory movements

Update purchase state

Create approval entry

Create audit log

COMMIT
```

### Collect Payment

```text
BEGIN

Lock outstanding invoices

Create payment

Create payment allocations

Update invoice paid/outstanding

Create audit log

COMMIT
```

These are the three most important DB transactions in the application.

---

## 43. PostgreSQL Isolation Level

Recommended initial isolation level:

```text
READ COMMITTED
```

Use explicit row locking where concurrency matters:

```sql
FOR UPDATE
```

Examples:

```text
Inventory
Outstanding invoices
Purchase approval
Document counters
```

There is likely no need for `SERIALIZABLE` isolation across the whole application.

---

## 44. Deadlock Prevention

For multi-product sales, do not lock product rows in arbitrary order.

Bad:

```text
Product C
Product A
Product B
```

Instead sort first:

```text
Product A
Product B
Product C
```

Then lock consistently.

The same rule applies when allocating one payment to multiple invoices.

---

## 45. Money Calculation Rule

Never use Go:

```go
float32
float64
```

for monetary business logic.

The database uses:

```text
NUMERIC
```

In Go, use either:

```text
integer paise
```

or a decimal library.

Example:

```text
₹125.50

12550 paise
```

However, inventory cost may require four decimal places, so using a decimal representation consistently across the domain may be cleaner.

---

## 46. Important Database Invariants

These should be enforced using a combination of PostgreSQL constraints and transactional Go code:

```text
Customer phone is unique.

Product SKU is unique.

Inventory is unique per store/product.

Inventory can never become negative.

Purchase approval happens once.

Approved purchase cannot be edited.

Finalized invoice cannot be directly edited.

Payment amount must be positive.

Payment allocations cannot exceed payment amount.

Payment allocation cannot exceed invoice outstanding.

Invoice paid amount cannot exceed invoice total.

Credit invoice must have a due date.

Every inventory change has an InventoryMovement.

Every collected payment identifies collected_by.

Every financial transaction remains auditable.

Automatic reminder stage cannot be duplicated.

Store Admin cannot operate outside assigned stores.
```

---

## 47. Tables That Should Not Be Primary Sources of Truth

### Avoid a primary `customer_credit` table

Do not make this authoritative:

```text
customer_id
total_due
```

Instead, derive authoritative credit from:

```text
SUM(invoice.outstanding_amount)
```

A reporting cache such as:

```text
customer_credit_summary
```

can be added later if needed.

### Avoid a separate `revenue` table

Revenue should come from invoices.

### Avoid a primary `profit` table

Do not store monthly profit as authoritative data.

Profit is derived from:

```text
Invoice Revenue
-
InvoiceItem COGS
-
Expenses
```

Materialized summaries can be introduced later.

---

## 48. Most Important Tables

The architectural core is:

```text
customers

products

purchases
purchase_items

inventory
inventory_movements

invoices
invoice_items

payments
payment_allocations

expenses

notifications

audit_logs
```

Two particularly important relationships are:

```text
Purchase
→ PurchaseItem
→ InventoryMovement
→ Inventory
```

and:

```text
Invoice
→ PaymentAllocation
→ Payment
```

These flows determine whether stock, dues, revenue, and profit remain correct.

---

## 49. Recommended Migration Order

Before writing Go repository code, create the schema in this order:

```text
1. Master tables
   users
   stores
   customers
   products
   suppliers

2. RBAC

3. Procurement

4. Inventory

5. Sales

6. Payments

7. Expenses

8. Configuration

9. Notifications

10. Audit / Outbox / Idempotency
```

---

## 50. Recommended Next Step

The next useful artifact should be the actual PostgreSQL migration SQL containing:

```text
CREATE TABLE
FOREIGN KEY
CHECK
UNIQUE
CREATE INDEX
```

That will force resolution of the remaining implementation-level database decisions, especially:

```text
Foreign-key delete behavior

Exact status constraints

Document-number generation

Payment reversal semantics

Transaction-safe inventory updates

Payment-allocation constraints
```
