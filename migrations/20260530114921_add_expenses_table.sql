-- +goose Up
CREATE TABLE expenses
(
    id         UUID PRIMARY KEY        DEFAULT gen_random_uuid(),
    title      VARCHAR(100)   NOT NULL,
    amount     DECIMAL(10, 2) NOT NULL,
    category   VARCHAR(50)    NOT NULL,
    date       DATE           NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE expenses;