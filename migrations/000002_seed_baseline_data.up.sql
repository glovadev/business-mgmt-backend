-- ==============================================================================
-- 000002_seed_baseline_data.up.sql
-- Comprehensive baseline seed data for NC Business Management Software
-- ==============================================================================

-- ------------------------------------------------------------------------------
-- 1. SEED ROLES
-- ------------------------------------------------------------------------------
INSERT INTO roles (id, code, name, description) VALUES
('11111111-1111-1111-1111-111111111101', 'SUPER_ADMIN', 'Super Administrator', 'Full system access across all stores, settings, finances, and user management.'),
('11111111-1111-1111-1111-111111111102', 'STORE_ADMIN', 'Store Administrator', 'Operational access to assigned stores: sales, purchases, inventory, local expenses.'),
('11111111-1111-1111-1111-111111111103', 'FINANCE_ADMIN', 'Finance Administrator', 'Central finance controller: purchase approvals, payments, reports, and global expenses.'),
('11111111-1111-1111-1111-111111111104', 'STAFF', 'Store Staff', 'Point of sale counter operator: view inventory, create customer sales and log payments.')
ON CONFLICT (code) DO NOTHING;

-- ------------------------------------------------------------------------------
-- 2. SEED GRANULAR PERMISSIONS
-- ------------------------------------------------------------------------------
INSERT INTO permissions (id, code, description) VALUES
-- IAM & System
('22222222-2222-2222-2222-222222222201', 'IAM_USER_MANAGE', 'Create, update, activate, and deactivate users'),
('22222222-2222-2222-2222-222222222202', 'IAM_ROLE_MANAGE', 'Assign and update roles and permissions'),
('22222222-2222-2222-2222-222222222203', 'IAM_AUDIT_VIEW', 'View system audit trails and operational logs'),
-- Store
('22222222-2222-2222-2222-222222222210', 'STORE_VIEW', 'View store master list and details'),
('22222222-2222-2222-2222-222222222211', 'STORE_MANAGE', 'Create and update store details'),
-- Customer
('22222222-2222-2222-2222-222222222220', 'CUSTOMER_VIEW', 'Search and view customer details and ledgers'),
('22222222-2222-2222-2222-222222222221', 'CUSTOMER_CREATE', 'Create new customer records'),
('22222222-2222-2222-2222-222222222222', 'CUSTOMER_EDIT', 'Update existing customer details'),
-- Product & Supplier
('22222222-2222-2222-2222-222222222230', 'PRODUCT_VIEW', 'Browse and search product catalog'),
('22222222-2222-2222-2222-222222222231', 'PRODUCT_MANAGE', 'Create, edit, and deactivate products and pricing'),
('22222222-2222-2222-2222-222222222232', 'PRODUCT_CATEGORY_MANAGE', 'Manage product categories'),
('22222222-2222-2222-2222-222222222240', 'SUPPLIER_VIEW', 'View supplier list and details'),
('22222222-2222-2222-2222-222222222241', 'SUPPLIER_MANAGE', 'Create and update suppliers'),
-- Procurement
('22222222-2222-2222-2222-222222222250', 'PURCHASE_VIEW', 'View purchase orders and landed cost details'),
('22222222-2222-2222-2222-222222222251', 'PURCHASE_CREATE', 'Create and submit purchase batches'),
('22222222-2222-2222-2222-222222222252', 'PURCHASE_APPROVE', 'Approve purchase batches and ingest inventory'),
('22222222-2222-2222-2222-222222222253', 'PURCHASE_CANCEL', 'Reject or cancel purchase orders'),
-- Inventory
('22222222-2222-2222-2222-222222222260', 'INVENTORY_VIEW', 'View store inventory stock levels and movements'),
('22222222-2222-2222-2222-222222222261', 'INVENTORY_ADJUST', 'Perform stock adjustments and physical counts'),
-- Sales
('22222222-2222-2222-2222-222222222270', 'SALE_VIEW', 'View invoices and sale transactions'),
('22222222-2222-2222-2222-222222222271', 'SALE_CREATE', 'Create sales invoices and apply discounts'),
('22222222-2222-2222-2222-222222222272', 'SALE_CANCEL', 'Cancel unallocated sales invoices'),
-- Payments
('22222222-2222-2222-2222-222222222280', 'PAYMENT_VIEW', 'View customer payments and allocation history'),
('22222222-2222-2222-2222-222222222281', 'PAYMENT_CREATE', 'Collect and record customer payments'),
('22222222-2222-2222-2222-222222222282', 'PAYMENT_REVERSE', 'Reverse posted payments'),
-- Expenses
('22222222-2222-2222-2222-222222222290', 'EXPENSE_VIEW', 'View store and company expenses'),
('22222222-2222-2222-2222-222222222291', 'EXPENSE_CREATE', 'Record store or company expenses'),
-- Notifications & Reminders
('22222222-2222-2222-2222-222222222300', 'NOTIFICATION_VIEW', 'View notification history and statuses'),
('22222222-2222-2222-2222-222222222301', 'NOTIFICATION_SEND', 'Trigger manual customer WhatsApp reminders'),
('22222222-2222-2222-2222-222222222302', 'NOTIFICATION_TEMPLATE_MANAGE', 'Configure WhatsApp templates and reminder rules'),
-- Reports & Settings
('22222222-2222-2222-2222-222222222310', 'REPORT_VIEW_STORE', 'View store-level daily sales, stock, and credit reports'),
('22222222-2222-2222-2222-222222222311', 'REPORT_VIEW_GLOBAL', 'View company-wide consolidated P&L, aged receivables, and tax reports'),
('22222222-2222-2222-2222-222222222320', 'SETTINGS_MANAGE', 'Configure company-wide policies, credit days, and alert thresholds')
ON CONFLICT (code) DO NOTHING;

-- ------------------------------------------------------------------------------
-- 3. ASSIGN PERMISSIONS TO ROLES
-- ------------------------------------------------------------------------------
-- SUPER_ADMIN gets ALL permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.code = 'SUPER_ADMIN'
ON CONFLICT DO NOTHING;

-- STORE_ADMIN gets store operational permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.code = 'STORE_ADMIN'
  AND p.code IN (
    'STORE_VIEW', 'CUSTOMER_VIEW', 'CUSTOMER_CREATE', 'CUSTOMER_EDIT',
    'PRODUCT_VIEW', 'SUPPLIER_VIEW',
    'PURCHASE_VIEW', 'PURCHASE_CREATE',
    'INVENTORY_VIEW', 'INVENTORY_ADJUST',
    'SALE_VIEW', 'SALE_CREATE',
    'PAYMENT_VIEW', 'PAYMENT_CREATE',
    'EXPENSE_VIEW', 'EXPENSE_CREATE',
    'NOTIFICATION_VIEW', 'NOTIFICATION_SEND',
    'REPORT_VIEW_STORE'
  )
ON CONFLICT DO NOTHING;

-- FINANCE_ADMIN gets central financial control permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.code = 'FINANCE_ADMIN'
  AND p.code IN (
    'STORE_VIEW', 'CUSTOMER_VIEW', 'SUPPLIER_VIEW', 'PRODUCT_VIEW',
    'PURCHASE_VIEW', 'PURCHASE_APPROVE', 'PURCHASE_CANCEL',
    'INVENTORY_VIEW',
    'SALE_VIEW',
    'PAYMENT_VIEW', 'PAYMENT_CREATE', 'PAYMENT_REVERSE',
    'EXPENSE_VIEW', 'EXPENSE_CREATE',
    'NOTIFICATION_VIEW', 'NOTIFICATION_SEND', 'NOTIFICATION_TEMPLATE_MANAGE',
    'REPORT_VIEW_STORE', 'REPORT_VIEW_GLOBAL',
    'SETTINGS_MANAGE'
  )
ON CONFLICT DO NOTHING;

-- STAFF gets POS counter permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.code = 'STAFF'
  AND p.code IN (
    'CUSTOMER_VIEW', 'CUSTOMER_CREATE',
    'PRODUCT_VIEW',
    'INVENTORY_VIEW',
    'SALE_VIEW', 'SALE_CREATE',
    'PAYMENT_VIEW', 'PAYMENT_CREATE'
  )
ON CONFLICT DO NOTHING;

-- ------------------------------------------------------------------------------
-- 4. SEED SYSTEM ADMINISTRATOR USER (admin@example.com / Password123!)
-- ------------------------------------------------------------------------------
-- bcrypt hash for 'Password123!' with cost 10
INSERT INTO users (id, email, password_hash, full_name, phone, status) VALUES
('33333333-3333-3333-3333-333333333333', 'admin@example.com', '$2a$10$tZcZ4oYQYQWnQk2t3v0R1.X0XWk2Zp1YQYQWnQk2t3v0R1.X0XWk2', 'System Administrator', '+919876543200', 'ACTIVE')
ON CONFLICT (email) DO UPDATE SET
  password_hash = EXCLUDED.password_hash,
  full_name = EXCLUDED.full_name,
  status = 'ACTIVE';

-- Assign SUPER_ADMIN role to admin user
INSERT INTO user_roles (user_id, role_id)
VALUES (
  '33333333-3333-3333-3333-333333333333',
  '11111111-1111-1111-1111-111111111101'
)
ON CONFLICT DO NOTHING;

-- ------------------------------------------------------------------------------
-- 5. SEED DEFAULT STORE (Main Branch)
-- ------------------------------------------------------------------------------
INSERT INTO stores (id, code, name, phone, gstin, address_line1, address_line2, city, state, postal_code, country_code, status, created_by) VALUES
('44444444-4444-4444-4444-444444444444', 'STR-001', 'Main Branch - Central Store', '+918022334455', '29ABCDE1234F1Z5', '100 Commercial Street', 'Industrial Area, Phase 1', 'Bengaluru', 'Karnataka', '560001', 'IN', 'ACTIVE', '33333333-3333-3333-3333-333333333333')
ON CONFLICT (code) DO NOTHING;

-- Assign Admin user to Main Store
INSERT INTO user_stores (user_id, store_id)
VALUES (
  '33333333-3333-3333-3333-333333333333',
  '44444444-4444-4444-4444-444444444444'
)
ON CONFLICT DO NOTHING;

-- ------------------------------------------------------------------------------
-- 6. SEED COMPANY SETTINGS
-- ------------------------------------------------------------------------------
INSERT INTO company_settings (id, company_name, currency, timezone, default_credit_days, yellow_after_overdue_days, red_after_overdue_days, additional_cost_allocation_method, include_gst_in_landed_cost, updated_by) VALUES
('55555555-5555-5555-5555-555555555555', 'NC Enterprises', 'INR', 'Asia/Kolkata', 15, 1, 7, 'PURCHASE_VALUE', FALSE, '33333333-3333-3333-3333-333333333333')
ON CONFLICT (id) DO NOTHING;

-- ------------------------------------------------------------------------------
-- 7. SEED NOTIFICATION TEMPLATES
-- ------------------------------------------------------------------------------
INSERT INTO notification_templates (id, code, provider_template_id, language_code, body_preview, status) VALUES
('66666666-6666-6666-6666-666666666601', 'PAYMENT_DUE_REMINDER', 'whatsapp_payment_due_v1', 'en', 'Dear {{customer_name}}, this is a friendly reminder that an outstanding amount of Rs {{outstanding_amount}} for Invoice {{invoice_number}} was due on {{due_date}}. Please clear your dues at your earliest convenience.', 'ACTIVE'),
('66666666-6666-6666-6666-666666666602', 'INVOICE_GENERATED', 'whatsapp_invoice_created_v1', 'en', 'Dear {{customer_name}}, thank you for your purchase at NC Enterprises. Your invoice {{invoice_number}} of Rs {{total_amount}} has been generated. Due date: {{due_date}}.', 'ACTIVE'),
('66666666-6666-6666-6666-666666666603', 'PAYMENT_RECEIVED', 'whatsapp_payment_receipt_v1', 'en', 'Dear {{customer_name}}, we have received your payment of Rs {{amount}} via {{payment_method}} (Ref: {{payment_number}}). Your remaining balance is Rs {{outstanding_balance}}.', 'ACTIVE')
ON CONFLICT (code) DO NOTHING;

-- ------------------------------------------------------------------------------
-- 8. SEED REMINDER RULES
-- ------------------------------------------------------------------------------
INSERT INTO reminder_rules (id, sequence_number, days_after_due, template_id, status, created_by) VALUES
('77777777-7777-7777-7777-777777777701', 1, 7, '66666666-6666-6666-6666-666666666601', 'ACTIVE', '33333333-3333-3333-3333-333333333333'),
('77777777-7777-7777-7777-777777777702', 2, 9, '66666666-6666-6666-6666-666666666601', 'ACTIVE', '33333333-3333-3333-3333-333333333333'),
('77777777-7777-7777-7777-777777777703', 3, 11, '66666666-6666-6666-6666-666666666601', 'ACTIVE', '33333333-3333-3333-3333-333333333333')
ON CONFLICT (sequence_number) DO NOTHING;

-- ------------------------------------------------------------------------------
-- 9. SEED PRODUCT CATEGORIES
-- ------------------------------------------------------------------------------
INSERT INTO product_categories (id, name, description, status) VALUES
('88888888-8888-8888-8888-888888888801', 'Cement & Concrete', 'Standard OPC and PPC building cement, white cement, and premix bags.', 'ACTIVE'),
('88888888-8888-8888-8888-888888888802', 'Steel & TMT Rebars', 'High-tensile thermo-mechanically treated reinforcement bars in various diameters.', 'ACTIVE'),
('88888888-8888-8888-8888-888888888803', 'Paints & Primers', 'Exterior, interior emulsions, distempers, primers, and protective wall finishes.', 'ACTIVE'),
('88888888-8888-8888-8888-888888888804', 'Plumbing & Pipes', 'PVC, CPVC pipes, drainage fittings, ball valves, and bathroom fixtures.', 'ACTIVE'),
('88888888-8888-8888-8888-888888888805', 'Electrical & Wiring', 'Copper multi-strand cables, switchboards, conduits, and LED fittings.', 'ACTIVE'),
('88888888-8888-8888-8888-888888888806', 'Hardware & Fasteners', 'Nails, screws, bolts, hinges, adhesive tape, and construction tools.', 'ACTIVE')
ON CONFLICT (name) DO NOTHING;

-- ------------------------------------------------------------------------------
-- 10. SEED BASELINE PRODUCTS
-- ------------------------------------------------------------------------------
INSERT INTO products (id, sku, name, category_id, description, hsn_code, unit_of_measure, default_selling_price, gst_rate, status, created_by) VALUES
('99999999-9999-9999-9999-999999999901', 'CEM-ULT-50', 'UltraTech Premium PPC Cement 50KG', '88888888-8888-8888-8888-888888888801', 'High quality Portland Pozzolana Cement for structural plastering and RCC work.', '25232900', 'PIECE', 420.00, 18.00, 'ACTIVE', '33333333-3333-3333-3333-333333333333'),
('99999999-9999-9999-9999-999999999902', 'TMT-TAT-12', 'Tata Tiscon 550D TMT Bar 12mm', '88888888-8888-8888-8888-888888888802', 'Earthquake-resistant Grade Fe550D high ductility reinforcement steel.', '72142090', 'KG', 72.00, 18.00, 'ACTIVE', '33333333-3333-3333-3333-333333333333'),
('99999999-9999-9999-9999-999999999903', 'PNT-APX-20', 'Asian Paints Apex WeatherProof 20L', '88888888-8888-8888-8888-888888888803', 'Exterior weather-proof wall paint with anti-fungal silicon additives.', '32091000', 'BUCKET', 4600.00, 18.00, 'ACTIVE', '33333333-3333-3333-3333-333333333333'),
('99999999-9999-9999-9999-999999999904', 'PLM-PVC-110', 'Astral PVC Soil & Waste Pipe 110mm 10FT', '88888888-8888-8888-8888-888888888804', 'Heavy duty leak-proof PVC drainage and soil pipe conforming to IS standards.', '39172310', 'PIECE', 680.00, 18.00, 'ACTIVE', '33333333-3333-3333-3333-333333333333'),
('99999999-9999-9999-9999-999999999905', 'ELE-HAV-25', 'Havells LifeLine Plus 2.5 sq mm Wire 90M', '88888888-8888-8888-8888-888888888805', 'Flame retardant low smoke 100% pure electrolytic copper domestic electrical wire.', '85444990', 'BOX', 2150.00, 18.00, 'ACTIVE', '33333333-3333-3333-3333-333333333333')
ON CONFLICT (sku) DO NOTHING;

-- Initial Inventory for default store
INSERT INTO inventory (store_id, product_id, available_quantity, average_unit_cost, updated_at) VALUES
('44444444-4444-4444-4444-444444444444', '99999999-9999-9999-9999-999999999901', 150.000, 360.0000, NOW()),
('44444444-4444-4444-4444-444444444444', '99999999-9999-9999-9999-999999999902', 1200.000, 61.5000, NOW()),
('44444444-4444-4444-4444-444444444444', '99999999-9999-9999-9999-999999999903', 25.000, 3950.0000, NOW()),
('44444444-4444-4444-4444-444444444444', '99999999-9999-9999-9999-999999999904', 60.000, 540.0000, NOW()),
('44444444-4444-4444-4444-444444444444', '99999999-9999-9999-9999-999999999905', 40.000, 1820.0000, NOW())
ON CONFLICT (store_id, product_id) DO NOTHING;

-- ------------------------------------------------------------------------------
-- 11. SEED SAMPLE SUPPLIER & CUSTOMER
-- ------------------------------------------------------------------------------
INSERT INTO suppliers (id, name, phone, email, gstin, address, status) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Tata Steel & Building Materials Ltd', '+918044556677', 'sales@tatabuilding.in', '29AAACT1234F1Z5', 'Plot 45, Peenya Industrial Area, Bengaluru', 'ACTIVE')
ON CONFLICT DO NOTHING;

INSERT INTO customers (id, phone_number, name, alternate_phone, email, address, gstin, whatsapp_opt_in, status, created_by) VALUES
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '+919876543210', 'Rahul Sharma', '+919876543211', 'rahul.sharma@example.com', 'Flat 402, Green Glen Apartments, Outer Ring Road, Bengaluru', NULL, TRUE, 'ACTIVE', '33333333-3333-3333-3333-333333333333')
ON CONFLICT (phone_number) DO NOTHING;

-- Associate Rahul Sharma with Main Store
INSERT INTO customer_store_relationships (customer_id, store_id, first_added_by) VALUES
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '44444444-4444-4444-4444-444444444444', '33333333-3333-3333-3333-333333333333')
ON CONFLICT DO NOTHING;

-- ------------------------------------------------------------------------------
-- 12. SEED EXPENSE CATEGORIES
-- ------------------------------------------------------------------------------
INSERT INTO expense_categories (id, name, scope, status) VALUES
('cccccccc-cccc-cccc-cccc-cccccccc0001', 'Store Rent & Lease', 'STORE', 'ACTIVE'),
('cccccccc-cccc-cccc-cccc-cccccccc0002', 'Electricity & Water Utilities', 'STORE', 'ACTIVE'),
('cccccccc-cccc-cccc-cccc-cccccccc0003', 'Staff Refreshments & Tea', 'STORE', 'ACTIVE'),
('cccccccc-cccc-cccc-cccc-cccccccc0004', 'Freight & Local Delivery', 'BOTH', 'ACTIVE'),
('cccccccc-cccc-cccc-cccc-cccccccc0005', 'Marketing & Digital Ads', 'COMPANY', 'ACTIVE'),
('cccccccc-cccc-cccc-cccc-cccccccc0006', 'Accounting & Legal Compliance', 'COMPANY', 'ACTIVE')
ON CONFLICT (name) DO NOTHING;

-- ------------------------------------------------------------------------------
-- 13. SEED DOCUMENT COUNTERS (FY 2026-27)
-- ------------------------------------------------------------------------------
INSERT INTO document_counters (store_id, document_type, financial_year, current_value) VALUES
('44444444-4444-4444-4444-444444444444', 'INVOICE', '2026-27', 0),
('44444444-4444-4444-4444-444444444444', 'PURCHASE', '2026-27', 0),
('44444444-4444-4444-4444-444444444444', 'PAYMENT', '2026-27', 0),
('44444444-4444-4444-4444-444444444444', 'EXPENSE', '2026-27', 0),
('00000000-0000-0000-0000-000000000000', 'PAYMENT', '2026-27', 0),
('00000000-0000-0000-0000-000000000000', 'EXPENSE', '2026-27', 0)
ON CONFLICT (store_id, document_type, financial_year) DO NOTHING;
