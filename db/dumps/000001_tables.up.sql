BEGIN;

CREATE TABLE IF NOT EXISTS orders(
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(250),
    track_number VARCHAR(250),
    entry TEXT,
    locale TEXT,
    internal_signature TEXT,
    customer_id TEXT,
    delivery_service TEXT,
    shardkey TEXT,
    sm_id INT,
    date_created TIMESTAMP,
    oof_chard TEXT
);

CREATE TABLE IF NOT EXISTS delivery(
    id SERIAL PRIMARY KEY,
    name VARCHAR(250),
    phone VARCHAR(250),
    zip TEXT,
    city TEXT,
    address TEXT,
    region TEXT,
    email TEXT,
    order_id INT, 
    FOREIGN KEY (order_id) REFERENCES orders (id)
);

CREATE TABLE IF NOT EXISTS payment(
    id SERIAL PRIMARY KEY,
    transaction VARCHAR(250),
    request_id VARCHAR(250),
    currency VARCHAR(250),
    provider VARCHAR(250),
    amount INT,
    payment_dt INT,
    bank VARCHAR(250),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT,
    order_id INT, 
    FOREIGN KEY (order_id) REFERENCES orders (id)
);

CREATE TABLE IF NOT EXISTS items(
    id SERIAL PRIMARY KEY,
    chrt_id INT,
    track_number VARCHAR(250),
    price INT,
    rid TEXT,
    name TEXT,
    sale INT,
    size TEXT,
    total_price INT,
    nm_id INT,
    brand TEXT,
    status INT,
    order_id INT,
    FOREIGN KEY (order_id) REFERENCES orders (id)
);

COMMIT;








