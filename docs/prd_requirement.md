# Business Management Software — Product Requirements

## 1. Product Overview

The Business Management Software will be a centralized system for managing a business operating through one or multiple stores.

The platform should manage:

- Products and inventory
- Procurement / stock addition
- Customers
- Sales
- Credit sales and customer dues
- Payment collection
- Expenses
- WhatsApp reminders
- Store-level and company-level financial reporting
- User roles and permissions
- Audit history

The system should provide store-level operational control while giving Finance Admins and Super Admins consolidated visibility across all stores.

---

## 2. Primary Business Goals

The system should help the business:

1. Know the exact inventory available in each store.
2. Record every product sale.
3. Track cash and credit sales separately.
4. Know how much every customer owes.
5. Automatically follow up with customers regarding overdue payments.
6. Track customer payment behavior.
7. Track store and company expenses.
8. Calculate revenue, cost of goods sold, gross profit, and net profit.
9. Know which employee performed important financial actions.
10. Produce consolidated financial information across stores.

---

## 3. Actors

### Customer

A Customer is a person purchasing products from one or more stores.

Capabilities:

- Purchase one or multiple products in a transaction.
- Make full payment during purchase.
- Make partial payment and keep the remaining amount as credit.
- Purchase entirely on credit.
- Make subsequent payments against outstanding dues.
- Receive WhatsApp payment reminders.
- Have multiple sales/invoices and multiple outstanding dues.

A customer is identified primarily using a **WhatsApp-enabled mobile number**.

#### Recommended rule

The phone number should be unique **across the company**, not separately within each store.

This avoids this situation:

```text
Store A
Customer: +91-9999999999

Store B
Customer: +91-9999999999
```

creating two different customers.

Instead:

```text
Customer
 ├── Purchases from Store A
 └── Purchases from Store B
```

This gives Finance/Super Admin a consolidated customer credit history.

---

## 4. User Roles and Authorization

### Store Admin

A Store Admin operates a particular store.

A Store Admin can:

- Login using email/password.
- View their assigned store.
- View store inventory.
- Add new stock purchases.
- Submit stock additions for approval.
- Add customers.
- Search existing customers using phone number.
- Create sales.
- Sell multiple products in one sale.
- Accept payments using:
  - Cash
  - UPI
  - Bank transfer
- Make full-payment sales.
- Make partial-credit sales.
- Make full-credit sales.
- Collect outstanding customer payments.
- Record store expenses.
- View sales belonging to their store.
- View customers and dues related to their store.
- View stock history of their store.

Store Admin must **not** be able to approve their own stock addition.

---

## 5. Finance Admin

Finance Admin primarily handles financial visibility and collections.

Finance Admin can:

- Login using email/password.
- Access all stores.
- View sales across stores.
- View revenue across stores.
- View expenses across stores.
- View customer outstanding credit.
- View overdue customers.
- View customer payment history.
- Collect customer dues.
- Record payments received.
- Send manual payment reminders.
- View automatic reminder history.
- Filter financial information by:
  - Store
  - Customer
  - Date range
  - Payment status
  - Payment method
  - Credit status
- View store-wise financial summaries.
- View consolidated company financial summaries.

Finance Admin should generally **not modify inventory** unless specifically granted permission.

---

## 6. Super Admin

Super Admin has full access to the platform.

Super Admin can perform everything available to:

- Store Admin
- Finance Admin

Additionally, Super Admin can:

- Create stores.
- Edit store information.
- Activate/deactivate stores.
- Create Store Admins.
- Create Finance Admins.
- Assign Store Admins to stores.
- Activate/deactivate users.
- Approve stock additions.
- Reject stock additions.
- Directly add approved inventory.
- Configure system-level settings.
- Configure credit reminder rules.
- Configure customer risk-category rules.
- Add company-level expenses.
- View consolidated dashboards.
- View audit logs.
- View inventory across all stores.

---

## 7. Core Modules

For future HLD, the product should be separated logically into these domains:

```text
Authentication & Authorization

Store Management

User Management

Customer Management

Product Catalog

Procurement / Stock Purchase

Inventory Management

Sales / Billing

Credit / Accounts Receivable

Payments / Collections

Expense Management

Notification Management

Financial Reporting

Configuration

Audit Logging
```

These do not necessarily need to become separate microservices. They are initially **business domains/modules**.

---

## 8. Store Management

A Store should contain:

```text
Store ID
Store Name
Address
Phone Number
Status
Created At
Created By
```

Statuses:

```text
ACTIVE
INACTIVE
```

All transactional records must contain a `store_id` where applicable.

Examples:

```text
Sale
Inventory
Store Expense
Stock Purchase
Payment Collection
```

Company-level expenses will not have a store associated with them.

---

## 9. Product Management

A Product should have information such as:

```text
Product ID
Product Name
SKU / Product Code
Category
Description
Selling Price
GST Rate
Unit of Measurement
Status
```

Example:

```text
Product:
Asian Paint 20L

SKU:
PAINT-ASIAN-20L

Selling Price:
₹4,500

GST:
18%

Unit:
Bucket
```

Product definition should preferably be company-wide.

Inventory quantity will be store-specific.

Therefore:

```text
Product ≠ Inventory
```

Example:

```text
Product:
Asian Paint 20L

Inventory:
Store A → 15 units
Store B → 8 units
Store C → 22 units
```

---

## 10. Procurement / Stock Purchase

Instead of Store Admin simply saying:

> Add 100 products to inventory

the system should model a **Stock Purchase / Procurement Batch**.

Example:

```text
Purchase #PUR-1024

Supplier: ABC Distributor
Store: Barpeta Road Store
Purchase Date: 01-Oct-2026

Products:
Product A × 100
Product B × 50
Product C × 30

Purchase Amount: ₹80,000
GST: ₹14,400
Transport: ₹5,000
Loading: ₹1,000
Other Charges: ₹500
```

This structure solves the logistic/GST requirement.

---

## 11. Stock Purchase Workflow

Store Admin:

```text
Create Purchase
      ↓
Add Products
      ↓
Add Purchase Costs
      ↓
Submit
      ↓
PENDING_APPROVAL
```

Super Admin:

```text
PENDING_APPROVAL
   ↓
Approve / Reject
```

On approval:

```text
Purchase = APPROVED

Inventory increases
Cost basis gets updated
Inventory ledger entry created
Audit event created
```

On rejection:

```text
Purchase = REJECTED

Inventory remains unchanged
```

Possible states:

```text
DRAFT
PENDING_APPROVAL
APPROVED
REJECTED
CANCELLED
```

### Important rule

Inventory must **not increase when the Store Admin submits the purchase**.

It increases only after Super Admin approval.

For purchases created directly by Super Admin:

```text
Create Purchase
      ↓
Automatically Approved
      ↓
Inventory Updated
```

---

## 12. Bulk Cost / Landed Cost Calculation

Additional costs can include:

```text
Transportation
Travel
Loading
Unloading
Insurance
Handling
GST
Other charges
```

Example:

```text
Product A purchase = ₹60,000
Product B purchase = ₹30,000
Product C purchase = ₹10,000

Total purchase = ₹100,000

Transportation = ₹5,000
```

For MVP, additional costs should be allocated **proportionally based on purchase value**.

Therefore:

```text
Product A → 60% → ₹3,000
Product B → 30% → ₹1,500
Product C → 10% → ₹500
```

The effective or landed cost becomes:

```text
Product A = purchase cost + allocated additional cost
```

Alternative allocation models can eventually include:

```text
By quantity
By purchase value
By weight
Manual allocation
```

But **purchase-value allocation** is a good MVP default.

### GST decision

If GST Input Tax Credit can be claimed, GST may not necessarily form part of inventory cost.

Therefore keep these separate:

```text
Base Purchase Cost
GST
Additional Charges
Landed Cost
```

This should be finalized with the business/accountant.

---

## 13. Inventory Management

Inventory should be maintained per:

```text
Store + Product
```

Example:

```text
Store A
Product X
Available Quantity = 57
```

However, the system should not rely solely on:

```text
quantity = 57
```

There should also be an **Inventory Movement Ledger**.

Example:

| Event | Quantity |
|---|---:|
| Purchase approved | +100 |
| Sale #1001 | -10 |
| Sale #1002 | -20 |
| Adjustment | -3 |
| Current | 67 |

Inventory movements could include:

```text
PURCHASE
SALE
SALE_RETURN
PURCHASE_RETURN
ADJUSTMENT_IN
ADJUSTMENT_OUT
STORE_TRANSFER_IN
STORE_TRANSFER_OUT
```

Even if some features like transfer/return are postponed, the inventory model should allow them.

---

## 14. Inventory Validation

The system must prevent:

```text
Available stock = 5

Customer purchase = 10
```

unless negative inventory is explicitly allowed.

Recommended:

> Negative inventory should not be allowed.

Sale confirmation should atomically validate and deduct stock.

---

## 15. Sales

A Sale/Invoice can contain multiple products.

Example:

```text
Invoice #INV-10023

Customer: Rahul
Store: Store A

Product A
2 × ₹1,000

Product B
3 × ₹500

Subtotal: ₹3,500
GST: ₹630

Total: ₹4,130
```

The invoice should contain:

```text
Invoice ID
Invoice Number
Store
Customer
Products
Quantity
Unit Price
Discount
GST
Subtotal
Total Amount
Amount Paid
Outstanding Amount
Payment Status
Sale Date
Credit Due Date
Created By
Created At
```

---

## 16. Payment Types During Sale

The system should support:

### Fully Paid

```text
Invoice = ₹5,000
Paid = ₹5,000
Due = ₹0
```

### Full Credit

```text
Invoice = ₹5,000
Paid = ₹0
Due = ₹5,000
```

### Partial Credit

```text
Invoice = ₹5,000
Paid = ₹2,000
Due = ₹3,000
```

Supporting **partial credit** from the start is recommended.

---

## 17. Payment Status

Invoice payment status:

```text
UNPAID
PARTIALLY_PAID
PAID
```

Credit status should be separate:

```text
NOT_DUE
OVERDUE
SETTLED
```

Example:

```text
Invoice Total = ₹5,000
Paid = ₹2,000

Payment Status:
PARTIALLY_PAID

Due Date:
10-Oct-2026

Current Date:
15-Oct-2026

Credit Status:
OVERDUE
```

---

## 18. Customer Credit

Credit should belong to **individual invoices**, rather than keeping only a single customer balance.

Example:

```text
Customer Rahul

Invoice 101
₹5,000
Due 10-Oct

Invoice 102
₹8,000
Due 20-Oct

Invoice 103
₹7,000
Due 30-Oct
```

Total credit can be calculated:

```text
₹20,000
```

But each invoice remains independently traceable.

---

## 19. Due Date

For any sale containing credit, a due date must exist.

Possible approach:

```text
Sale Date:
1 Oct

Default Credit Period:
15 days

Due Date:
16 Oct
```

Super Admin should configure:

```text
Default credit period
```

Store Admin may optionally override the due date if authorized.

---

## 20. Customer Risk Classification

Customers should be categorized:

```text
GREEN
YELLOW
RED
```

Recommended definition:

Classification should depend on **the oldest unpaid overdue invoice**.

Example configuration:

```text
GREEN:
0 days overdue

YELLOW:
1–6 days overdue

RED:
7+ days overdue
```

Example:

```text
Invoice A:
2 days overdue

Invoice B:
10 days overdue
```

Customer classification:

```text
RED
```

because the oldest unpaid overdue invoice is 10 days overdue.

Super Admin should configure thresholds.

---

## 21. Automatic WhatsApp Reminder

Super Admin should configure reminder rules.

Example:

```text
Reminder 1 → 7 days overdue
Reminder 2 → 9 days overdue
Reminder 3 → 11 days overdue
```

System workflow:

```text
Scheduled Job
     ↓
Find overdue invoices
     ↓
Calculate overdue days
     ↓
Evaluate reminder rules
     ↓
Check whether reminder already sent
     ↓
Send WhatsApp message
     ↓
Store notification result
```

---

## 22. Notification History

Every WhatsApp attempt should be recorded.

Suggested fields:

```text
Notification ID
Customer
Invoice
Phone Number
Notification Type
Trigger Type
Overdue Days
Message Template
Sent At
Status
Provider Message ID
Error Message
Triggered By
```

Trigger type:

```text
AUTOMATIC
MANUAL
```

Status:

```text
PENDING
SENT
DELIVERED
FAILED
```

The system must prevent duplicate automatic reminders.

---

## 23. Manual Reminder

Finance Admin and Super Admin should be able to manually send reminders.

Potentially Store Admin can also be allowed, but that should be permission-controlled.

Example:

```text
Customer: Rahul
Outstanding: ₹12,000
Oldest Due: 10 days

[Send Reminder]
```

The reminder should create a Notification History record.

---

## 24. Customer Payment Collection

Customer may later pay through:

```text
Cash
UPI
Bank Transfer
```

Payment record should contain:

```text
Payment ID
Customer
Store
Invoice
Amount
Payment Method
Transaction Reference
Collected By
Collected At
Notes
Created At
```

`collected_by` must identify the user.

Example:

```text
Payment:
₹5,000

Collected By:
Finance Admin / Elvis

Source:
UPI
```

Every payment must record the employee responsible for collecting/recording it.

---

## 25. Payment Allocation

Example:

```text
Invoice A due = ₹5,000
Invoice B due = ₹7,000

Customer pays ₹6,000
```

Recommended default:

**Oldest outstanding invoice first.**

Therefore:

```text
Invoice A:
₹5,000 → Settled

Invoice B:
₹1,000 paid
₹6,000 remaining
```

Finance Admin should be able to explicitly select an invoice if required.

---

## 26. Expense Management

Two kinds of expenses should exist.

### Store Expenses

Examples:

```text
Tea / Coffee
Electricity
Store Rent
Staff Travel
Maintenance
Internet
Miscellaneous
```

Store Admin can create these.

### Company Expenses

Examples:

```text
Head Office Rent
Company Software
CA Fees
Company Marketing
Legal Expenses
Corporate Travel
```

Only Super Admin should add company-level expenses unless permission is delegated.

Expense should contain:

```text
Expense ID
Store ID (nullable)
Expense Category
Amount
Expense Date
Payment Method
Vendor
Description
Receipt / Attachment
Created By
Created At
```

If `store_id = NULL`, it represents a company-level expense.

---

## 27. Financial Definitions

These definitions should be explicit because terms like revenue and collections are often confused.

### Gross Sales

Total invoice amount generated.

```text
Gross Sales = Σ invoices
```

### Net Sales / Revenue

```text
Gross Sales
- discounts
- sales returns
```

### Cash Collected

Actual money received.

Example:

```text
Sale = ₹10,000
Customer paid = ₹3,000
Credit = ₹7,000

Revenue = ₹10,000
Cash collected = ₹3,000
Receivable = ₹7,000
```

### Cost of Goods Sold

```text
COGS = quantity sold × inventory cost
```

### Gross Profit

```text
Revenue - COGS
```

### Net Profit

```text
Gross Profit
- Store Expenses
- Company Expenses
```

---

## 28. Product Costing Method

Products may be purchased at different prices.

Example:

```text
Jan:
100 units × ₹100

Feb:
100 units × ₹120
```

For MVP, use:

### Weighted Average Cost

After January:

```text
Average cost = ₹100
```

After February:

```text
(100 × 100 + 100 × 120) / 200

= ₹110
```

Then if one product is sold for ₹150:

```text
Revenue = ₹150
COGS = ₹110
Gross Profit = ₹40
```

This is considerably simpler than FIFO lot accounting while still supporting meaningful profit reporting.

---

## 29. Dashboard

### Store Admin Dashboard

Suggested information:

```text
Today's Sales
Today's Collection
Current Outstanding Credit
Today's Expenses
Low Stock Products
Overdue Customers
Inventory Value
```

Only assigned store.

### Finance Admin Dashboard

Suggested:

```text
Total Sales
Total Collections
Outstanding Receivables
Overdue Amount
Store-wise Revenue
Store-wise Expenses
Customer Credit Aging
Payments Collected
```

Filters:

```text
Today
Yesterday
This Week
This Month
Custom Range
Store
```

### Super Admin Dashboard

Suggested:

```text
Total Sales
Revenue
Cash Collected
Outstanding Credit
Overdue Amount
COGS
Gross Profit
Store Expenses
Company Expenses
Net Profit
Inventory Value
Store-wise Performance
Red/Yellow Customers
Low Stock Alerts
```

---

## 30. Credit Aging Report

Example:

| Due Period | Outstanding |
|---|---:|
| Not Due | ₹2,00,000 |
| 1–7 days overdue | ₹50,000 |
| 8–30 days | ₹30,000 |
| 31–60 days | ₹15,000 |
| 60+ days | ₹10,000 |

Finance Admin should be able to drill down into the customers/invoices contributing to these numbers.

---

## 31. Audit Logging

Audit logging should be treated as a **core requirement**.

Important actions should record:

```text
Actor
Action
Entity
Previous Value
New Value
Timestamp
Store
```

Example:

```text
Actor:
StoreAdmin-123

Action:
CREATE_PAYMENT

Invoice:
INV-10023

Amount:
₹5,000

Timestamp:
02-Oct-2026 14:31
```

Important audited actions:

```text
Login

Create/update customer

Create sale

Cancel sale

Create payment

Modify payment

Create expense

Update expense

Submit stock

Approve stock

Reject stock

Inventory adjustment

Send manual reminder

Change reminder rules

Create admin

Disable admin
```

---

## 32. Editing Financial Transactions

Finalized transactions should not be freely edited.

Instead use controlled actions such as:

```text
Cancel Sale
Refund
Credit Note
Inventory Reversal
Payment Reversal
```

and maintain full history.

The same principle should apply to payments.

---

## 33. Authentication

Admin users:

```text
Store Admin
Finance Admin
Super Admin
```

authenticate using:

```text
Email
Password
```

Recommended later:

```text
Forgot Password
Password Reset
Session Expiration
Account Disable
Optional 2FA for Super Admin
```

Passwords must never be stored directly; only secure password hashes should be stored.

---

## 34. Authorization

The system should implement Role-Based Access Control.

Conceptually:

```text
USER
 ↓
ROLE
 ↓
PERMISSIONS
 ↓
STORE SCOPE
```

Example:

```text
Store Admin

Permissions:
CREATE_SALE
CREATE_CUSTOMER
CREATE_PAYMENT
CREATE_EXPENSE
CREATE_STOCK_REQUEST

Scope:
Store A
```

Finance Admin:

```text
VIEW_SALES
VIEW_EXPENSES
VIEW_RECEIVABLES
CREATE_PAYMENT
SEND_REMINDER

Scope:
ALL STORES
```

This is cleaner than scattering role-specific checks throughout the application.

---

## 35. Important Business Invariants

These rules should eventually translate directly into application/database validation.

```text
Customer WhatsApp phone must be unique.

Product quantity sold cannot exceed available inventory.

Inventory only increases after purchase approval.

Every inventory change must create an inventory movement.

Every payment must identify who collected/recorded it.

Outstanding amount cannot become negative.

Paid amount cannot exceed invoice amount unless advance payments are explicitly supported.

Credit sales must have a due date.

A finalized invoice should not be directly modified.

Automatic reminder for the same reminder stage must not be duplicated.

Store Admin cannot approve their own stock request.

Users cannot access stores outside their assigned scope.

Every financial transaction must be auditable.
```

---

## 36. Suggested High-Level Data Model

Core entities:

```text
User
Role
Permission
UserStore

Store

Customer
CustomerStoreRelationship

Product
ProductCategory

Supplier

Purchase
PurchaseItem
PurchaseAdditionalCost
PurchaseApproval

Inventory
InventoryMovement

Sale / Invoice
SaleItem

Payment
PaymentAllocation

Expense
ExpenseCategory

CreditConfiguration
ReminderConfiguration
CustomerRiskConfiguration

Notification
NotificationTemplate

AuditLog
```

Simplified relationship:

```text
                     Company
                        |
           +------------+------------+
           |                         |
        Store A                   Store B
           |                         |
       Inventory                 Inventory
           |                         |
         Sales                     Sales
           \                         /
            \                       /
                 Customer
                    |
                 Invoice
                    |
              Outstanding
                    |
                 Payment
```

---

## 37. System Jobs

### Reminder Scheduler

```text
Runs periodically
→ finds overdue invoices
→ evaluates reminder configuration
→ sends WhatsApp messages
```

### Customer Risk Evaluator

```text
Find outstanding invoices
→ calculate maximum overdue days
→ assign Green / Yellow / Red
```

This calculation can happen dynamically or periodically.

### Financial Aggregation

Potentially calculates/reporting aggregates such as:

```text
Daily sales
Monthly revenue
Store totals
Outstanding amounts
```

Whether aggregation is precomputed or calculated dynamically should be decided during HLD based on volume.

---

## 38. Error / Edge Cases

### Concurrent Sales

Two Store Admins attempt to sell the last unit simultaneously.

The system must ensure that only one succeeds.

### Payment Retry

Admin submits payment and network connection drops.

Retry must not accidentally create two payments.

### WhatsApp Failure

WhatsApp provider may:

```text
Timeout
Reject message
Return failure
```

The system should record the failure and allow retry.

### Purchase Approval Conflict

A Store Admin should not be able to modify a purchase while Super Admin is approving it.

### Customer Already Exists

Store Admin enters an existing phone number.

System should show the existing customer instead of creating another one.

### Incorrect Payment

The system should provide a controlled reversal instead of deleting the payment.

### Stock Correction

Physical inventory may differ from system inventory.

Authorized users should have:

```text
Inventory Adjustment
```

with mandatory reason and audit history.

---

## 39. Non-Functional Requirements

### Security

- HTTPS everywhere.
- Password hashing.
- Authorization enforced server-side.
- Secure session/JWT management.
- Rate limiting on authentication APIs.
- Sensitive data should not appear unnecessarily in logs.

### Data Consistency

Financial and inventory operations should be transactional.

Example:

```text
Create Sale
+
Deduct Inventory
+
Create Receivable
```

should either all succeed or all fail.

### Idempotency

Critical operations such as:

```text
Payment creation
Sale creation
WhatsApp notification
```

should support duplicate-request protection.

### Auditability

Important financial changes must remain traceable.

### Backup

Database backups should be performed automatically.

### Observability

The application should provide:

```text
Application logs
Error logs
API metrics
Background-job monitoring
WhatsApp failure monitoring
```

### Performance

Common operations such as:

```text
Create sale
Search customer
Load inventory
View dashboard
```

should remain responsive as transaction count increases.

---

## 40. Suggested MVP Scope

### V1

```text
Authentication + RBAC

Store management

Admin management

Customer management

Product catalog

Stock purchase + approval

Inventory management

Sales

Full / partial credit

Payment collection

Expenses

Credit status

WhatsApp reminders

Audit logs

Basic dashboards and reports
```

### Phase 2

```text
Supplier payments/accounts payable

Purchase returns

Sales returns/refunds

Inter-store inventory transfer

Barcode scanning

Invoice PDF generation

GST-compliant accounting reports

Advanced accounting integration

Employee payroll

Credit limits

Customer advance balance

Mobile app

Multi-company SaaS
```

---

## 41. Important Decisions Before HLD

1. **Can Store Admin belong to only one store or multiple stores?**  
   Recommendation: support multiple technically, even if MVP assigns one.

2. **Can the same customer purchase from multiple stores?**  
   Recommendation: yes; keep customer company-wide.

3. **Can sales be partially paid?**  
   Recommendation: yes.

4. **How is a customer payment allocated among multiple invoices?**  
   Recommendation: explicit invoice selection with oldest-first default.

5. **What costing method determines COGS?**  
   Recommendation: weighted-average costing.

6. **How should transport/other procurement costs be allocated?**  
   Recommendation: purchase-value proportion for MVP.

7. **Does GST form part of product cost?**  
   Needs accounting/business clarification based on GST/ITC treatment.

8. **Can Store Admin override credit due dates?**  
   Needs business decision.

9. **Can sales/payments be deleted?**  
   Recommendation: no; reverse/cancel instead.

10. **How exactly is Green/Yellow/Red calculated?**  
    Recommendation: oldest outstanding overdue invoice.

---

## 42. Most Important Conceptual Separation

When moving to HLD/LLD, these should remain separate concepts:

```text
Product
    ≠
Inventory

Purchase
    ≠
Inventory

Sale
    ≠
Payment

Revenue
    ≠
Cash Collection

Credit
    ≠
Payment

Expense
    ≠
COGS

Customer
    ≠
Store Customer

Notification
    ≠
Credit

Inventory Quantity
    ≠
Inventory Movement
```

---

## 43. Recommended Design Sequence

```text
PRD
 ↓
Business Rules / Edge Cases
 ↓
Core Entities + Relationships
 ↓
DB Schema / ER Diagram
 ↓
Major API Contracts
 ↓
HLD
 ↓
Detailed Flows / Sequence Diagrams
 ↓
LLD
```

The recommended next step is to create the **domain model + ER diagram**, especially around:

```text
Customer → Sale → Payment → Credit
```

and:

```text
Purchase → Inventory → Sale → COGS
```

Once those relationships are correct, the HLD becomes much easier and data-model issues can be identified before implementation.
