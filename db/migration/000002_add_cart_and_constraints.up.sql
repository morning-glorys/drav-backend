ALTER TABLE sellers ADD CONSTRAINT unique_user_seller UNIQUE (user_id);


CREATE TABLE carts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_cart UNIQUE (user_id)
);


CREATE TABLE cart_items (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1,
    CONSTRAINT unique_cart_product UNIQUE (cart_id, product_id)
);