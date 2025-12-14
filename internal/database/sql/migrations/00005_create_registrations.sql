-- +goose Up
-- +goose StatementBegin
CREATE TYPE payment_status AS ENUM ('pending', 'confirmed', 'cancelled', 'refunded');

CREATE TABLE registrations (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ticket_batch_id BIGINT NOT NULL REFERENCES ticket_batches(id) ON DELETE RESTRICT,
    amount_paid DECIMAL(10, 2) NOT NULL CHECK (amount_paid >= 0),
    payment_status payment_status DEFAULT 'pending',
    registration_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    payment_confirmed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    UNIQUE(event_id, user_id)
);

CREATE INDEX idx_registrations_event_id ON registrations(event_id);
CREATE INDEX idx_registrations_user_id ON registrations(user_id);
CREATE INDEX idx_registrations_ticket_batch_id ON registrations(ticket_batch_id);
CREATE INDEX idx_registrations_payment_status ON registrations(payment_status);
CREATE INDEX idx_registrations_registration_date ON registrations(registration_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS registrations;
DROP TYPE IF EXISTS payment_status;
-- +goose StatementEnd