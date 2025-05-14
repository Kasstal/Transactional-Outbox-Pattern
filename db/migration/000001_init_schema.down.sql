-- Down migration (rollback) script
-- Executes in reverse order of creation to respect dependencies

-- 1. First drop constraints
ALTER TABLE IF EXISTS history
DROP CONSTRAINT IF EXISTS history_order_id_fkey;

ALTER TABLE IF EXISTS order_items
DROP CONSTRAINT IF EXISTS order_items_order_id_fkey;

ALTER TABLE IF EXISTS payments
DROP CONSTRAINT IF EXISTS payments_order_id_fkey;

-- 2. Drop indexes
DROP INDEX IF EXISTS payments_type_idx;
DROP INDEX IF EXISTS payments_order_id_idx;
DROP INDEX IF EXISTS order_items_status_idx;
DROP INDEX IF EXISTS order_items_product_id_idx;
DROP INDEX IF EXISTS order_items_order_id_idx;
DROP INDEX IF EXISTS history_date_idx;
DROP INDEX IF EXISTS history_order_id_idx;
DROP INDEX IF EXISTS orders_general_id_idx;
DROP INDEX IF EXISTS orders_created_at_idx;
DROP INDEX IF EXISTS orders_platform_idx;
DROP INDEX IF EXISTS orders_status_idx;
DROP INDEX IF EXISTS outbox_events_created_at_idx;
DROP INDEX IF EXISTS outbox_events_status_idx;
DROP INDEX IF EXISTS outbox_events_aggregate_idx;

-- 3. Drop tables in reverse dependency order
DROP TABLE IF EXISTS outbox_events;
DROP TABLE IF EXISTS history;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS inbox_events;
-- 4. Drop custom types
DROP TYPE IF EXISTS payment_type;

-- 5. Drop extensions
DROP EXTENSION IF EXISTS pgcrypto;