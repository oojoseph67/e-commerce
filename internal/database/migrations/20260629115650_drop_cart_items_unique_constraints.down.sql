-- Recreate the unique constraints (rollback)
ALTER TABLE cart_items ADD CONSTRAINT cart_items_cart_id_key UNIQUE (cart_id);
ALTER TABLE cart_items ADD CONSTRAINT cart_items_product_id_key UNIQUE (product_id);