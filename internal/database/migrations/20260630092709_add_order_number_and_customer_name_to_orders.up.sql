-- Add order_number (human-readable order ID) and customer_name to orders table
ALTER TABLE orders ADD COLUMN IF NOT EXISTS order_number VARCHAR(50) UNIQUE NOT NULL DEFAULT 'TEMP-' || gen_random_uuid();
ALTER TABLE orders ADD COLUMN IF NOT EXISTS customer_name VARCHAR(255) NOT NULL DEFAULT '';

-- Create index on order_number for fast lookups
CREATE INDEX IF NOT EXISTS idx_orders_order_number ON orders(order_number);