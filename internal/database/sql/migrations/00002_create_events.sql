-- +goose Up
-- +goose StatementBegin
CREATE TYPE event_status AS ENUM ('draft', 'published', 'cancelled', 'finished');

CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    city VARCHAR(50),
    street VARCHAR(50),
    neighbourhood VARCHAR(50),
    max_capacity INTEGER,
    status event_status DEFAULT 'draft',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_dates CHECK (end_date >= start_date)
);

CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_start_date ON events(start_date);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
DROP TYPE IF EXISTS event_status;
-- +goose StatementEnd