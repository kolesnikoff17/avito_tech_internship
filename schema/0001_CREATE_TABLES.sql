CREATE TABLE users (
    user_id INTEGER CHECK (user_id > 0) PRIMARY KEY,
    amount DECIMAL(18,2) CHECK ( amount >= 0 ) NOT NULL,
    reserved DECIMAL(18,2) CHECK ( reserved >= 0 ) NOT NULL DEFAULT 0.00
);

CREATE TABLE status (
    status_id SERIAL PRIMARY KEY,
    status_name VARCHAR(20)
);

CREATE TABLE services (
    service_id SERIAL PRIMARY KEY,
    service_name VARCHAR(55)
);

CREATE TABLE replenishments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    amount DECIMAL(18,2) CHECK ( amount >= 0 ),
    created TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users (user_id)
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    order_id INTEGER UNIQUE,
    service_id INTEGER,
    user_id INTEGER,
    order_sum DECIMAL(18,2) CHECK ( order_sum >= 0 ),
    status_id INTEGER,
    created TIMESTAMPTZ NOT NULL DEFAULT now(),
    modified TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users (user_id),
    FOREIGN KEY (status_id) REFERENCES status (status_id),
    FOREIGN KEY (service_id) REFERENCES services (service_id)
);

INSERT INTO services (service_id, service_name) VALUES
    (1, 'Rent'),
    (2, 'Good bought'),
    (3, 'Advertisement bought'),
    (4, 'Service commission'),
    (5, 'Monthly subscription payment');

INSERT INTO status (status_id, status_name) VALUES
    (1, 'Pending'),
    (2, 'Approved'),
    (3, 'Canceled');