-- +goose Up
-- +goose StatementBegin
CREATE TABLE ticket_batches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    batch_name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    date_limit TIMESTAMPTZ NOT NULL,
    quantity_limit INTEGER NOT NULL,
    batch_order INTEGER NOT NULL,
    active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    CONSTRAINT check_at_least_one_limit CHECK (
        date_limit IS NOT NULL OR quantity_limit IS NOT NULL
    )
);

CREATE INDEX idx_ticket_batches_event_id ON ticket_batches(event_id);
CREATE INDEX idx_ticket_batches_active ON ticket_batches(active);
CREATE INDEX idx_ticket_batches_order ON ticket_batches(batch_order);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ticket_batches;
-- +goose StatementEnd