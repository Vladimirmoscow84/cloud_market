BEGIN;

CREATE TABLE IF NOT EXISTS orders(
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(250) NOT NULL UNIQUE,
    track_number VARCHAR(250) NOT NULL UNIQUE,
    entry TEXT,
    locale TEXT NOT NULL,
    internal_signature TEXT,
    customer_id TEXT NOT NULL,
    delivery_service TEXT NOT NULL,
    shardkey TEXT NOT NULL,
    sm_id INT,
    date_created TEXT NOT NULL,
    oof_chard TEXT NOT NULL 
);

CREATE TABLE IF NOT EXISTS delivery(
    id SERIAL PRIMARY KEY,
    name VARCHAR(250) NOT NULL,
    phone VARCHAR(250) NOT NULL,
    zip TEXT NOT NULL,
    city TEXT NOT NULL,
    address TEXT NOT NULL,
    region TEXT NOT NULL,
    email TEXT NOT NULL,
    order_id INT, 
    FOREIGN KEY (order_id) REFERENCES orders (id)
);

CREATE TABLE IF NOT EXISTS payment(
    id SERIAL PRIMARY KEY,
    transaction VARCHAR(250) NOT NULL UNIQUE,
    request_id VARCHAR(250),
    currency VARCHAR(250) NOT NULL,
    provider VARCHAR(250) NOT NULL,
    amount INT,
    payment_dt INT,
    bank VARCHAR(250) NOT NULL,
    delivery_cost INT,
    goods_total INT,
    custom_fee INT,
    order_id INT, 
    FOREIGN KEY (order_id) REFERENCES orders (id)
);

CREATE TABLE IF NOT EXISTS items(
    id SERIAL PRIMARY KEY,
    chrt_id INT,
    track_number VARCHAR(250) NOT NULL UNIQUE,
    price INT,
    rid VARCHAR(250) NOT NULL UNIQUE,
    name VARCHAR(250) NOT NULL,
    sale INT,
    size VARCHAR(250) NOT NULL,
    total_price INT,
    nm_id INT,
    brand VARCHAR(250),
    status INT,
    order_id INT,
    FOREIGN KEY (order_id) REFERENCES orders (id)
);

COMMIT;








