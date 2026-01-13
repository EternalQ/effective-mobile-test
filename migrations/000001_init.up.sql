CREATE TABLE subsriptions (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR NOT NULL,
    price INT NOT NULL CHECK (price > 0),
    user_id UUID NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP
);

CREATE INDEX IF NOT EXISTS subsriptions_calc_idx ON subsriptions (user_id, service_name);