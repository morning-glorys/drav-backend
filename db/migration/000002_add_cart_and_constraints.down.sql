ALTER TABLE cart_items DROP CONSTRAINT IF EXISTS unique_cart_product;
ALTER TABLE carts DROP CONSTRAINT IF EXISTS unique_user_cart;
ALTER TABLE sellers DROP CONSTRAINT IF EXISTS unique_user_seller;
