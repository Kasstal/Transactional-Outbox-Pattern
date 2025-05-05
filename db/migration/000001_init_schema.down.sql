-- Drop foreign key constraints first (reverse order of creation)
ALTER TABLE card_payment_data DROP CONSTRAINT IF EXISTS card_payment_data_payment_id_fkey;
ALTER TABLE credit_data DROP CONSTRAINT IF EXISTS credit_data_payment_id_fkey;
ALTER TABLE payments DROP CONSTRAINT IF EXISTS payments_order_id_fkey;
ALTER TABLE order_items DROP CONSTRAINT IF EXISTS order_items_order_id_fkey;
ALTER TABLE history DROP CONSTRAINT IF EXISTS history_order_id_fkey;

-- Drop indexes (reverse order of creation)
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

-- Drop tables (reverse order of creation)
DROP TABLE IF EXISTS card_payment_data;
DROP TABLE IF EXISTS credit_data;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS history;
DROP TABLE IF EXISTS orders;

-- drop the enum type
DROP TYPE IF EXISTS payment_type;

-- Drop extension if needed
DROP EXTENSION IF EXISTS pgcrypto;