-- Remove order_number and customer_name from orders table
DROP INDEX IF EXISTS idx_orders_order_number;
ALTER TABLE orders DROP COLUMN IF EXISTS customer_name;
ALTER TABLE orders DROP COLUMN IF EXISTS order_number;