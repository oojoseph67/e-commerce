-- Rename price to subtotal for semantic clarity
-- Price stored the line total (qty * unit price), not the unit price
ALTER TABLE order_items RENAME COLUMN price TO subtotal;