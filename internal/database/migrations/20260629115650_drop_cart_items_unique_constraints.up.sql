-- Drop the incorrect unique constraints on cart_items
-- These constraints prevented a cart from having multiple items
-- and prevented a product from appearing in multiple carts

ALTER TABLE cart_items DROP CONSTRAINT IF EXISTS cart_items_cart_id_key;
ALTER TABLE cart_items DROP CONSTRAINT IF EXISTS cart_items_product_id_key;