CREATE TABLE IF NOT EXISTS cart_items (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	cart_id UUID UNIQUE NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
	product_id UUID UNIQUE NOT NULL REFERENCES products(id) ON DELETE CASCADE,
	quantity INTEGER NOT NULL CHECK (quantity > 0),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP WITH TIME ZONE
	-- UNIQUE(cart_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_cart_items_cart_id ON cart_items(cart_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_product_id ON cart_items(product_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_deleted_at ON cart_items(deleted_at);
