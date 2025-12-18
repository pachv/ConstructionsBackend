-- orders
CREATE TABLE orders (
    id VARCHAR(255),
    user_id VARCHAR(255),

    company_name VARCHAR(255),
    email VARCHAR(255),
    delivery_address VARCHAR(1024),
    comment TEXT,

    customer_name VARCHAR(255),
    customer_phone VARCHAR(255),
    consent BOOLEAN,

    created_at TIMESTAMP
);

-- order items
CREATE TABLE order_items (
    id VARCHAR(255),
    order_id VARCHAR(255),

    product_id VARCHAR(255),
    qty INT,

    -- снапшот чтобы в будущем заказ не ломался, если товар поменяется
    product_title VARCHAR(255),
    product_price INT,
    product_old_price INT,
    product_sale_percent INT,
    product_image_path VARCHAR(1024),

    created_at TIMESTAMP
);
