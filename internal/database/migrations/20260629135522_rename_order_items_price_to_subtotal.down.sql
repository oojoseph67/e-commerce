-- Rollback: rename subtotal back to price
ALTER TABLE order_items RENAME COLUMN subtotal TO price;