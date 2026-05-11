ALTER TABLE sellers
ADD CONSTRAINT IF NOT EXISTS unique_user_seller UNIQUE (user_id);

ALTER TABLE carts
ADD CONSTRAINT IF NOT EXISTS unique_user_cart UNIQUE (user_id);

ALTER TABLE cart_items
ADD CONSTRAINT IF NOT EXISTS unique_cart_product UNIQUE (cart_id, product_id);
