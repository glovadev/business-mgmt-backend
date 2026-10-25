-- ==============================================================================
-- 000001_init_schema.down.sql
-- Reverse 000001_init_schema.up.sql
-- ==============================================================================

DROP TABLE IF EXISTS document_counters CASCADE;
DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS outbox_events CASCADE;
DROP TABLE IF EXISTS idempotency_keys CASCADE;

DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS reminder_rules CASCADE;
DROP TABLE IF EXISTS notification_templates CASCADE;
DROP TABLE IF EXISTS company_settings CASCADE;

DROP TABLE IF EXISTS expenses CASCADE;
DROP TABLE IF EXISTS expense_categories CASCADE;

DROP TABLE IF EXISTS payment_allocations CASCADE;
DROP TABLE IF EXISTS payments CASCADE;

DROP TABLE IF EXISTS invoice_items CASCADE;
DROP TABLE IF EXISTS invoices CASCADE;

DROP TABLE IF EXISTS inventory_movements CASCADE;
DROP TABLE IF EXISTS inventory CASCADE;

DROP TABLE IF EXISTS purchase_approvals CASCADE;
DROP TABLE IF EXISTS purchase_additional_costs CASCADE;
DROP TABLE IF EXISTS purchase_items CASCADE;
DROP TABLE IF EXISTS purchases CASCADE;

DROP TABLE IF EXISTS suppliers CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS product_categories CASCADE;

DROP TABLE IF EXISTS customer_store_relationships CASCADE;
DROP TABLE IF EXISTS customers CASCADE;

DROP TABLE IF EXISTS user_sessions CASCADE;
DROP TABLE IF EXISTS user_stores CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

-- Drop store FK to user before dropping users and stores
ALTER TABLE IF EXISTS stores DROP CONSTRAINT IF EXISTS fk_stores_created_by;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS stores CASCADE;
