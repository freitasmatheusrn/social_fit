-- +goose Up
-- +goose StatementBegin
CREATE TYPE organizer_role AS ENUM ('main', 'collaborator');

CREATE TABLE event_organizers (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role organizer_role DEFAULT 'collaborator',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(event_id, user_id)
);

CREATE INDEX idx_event_organizers_event_id ON event_organizers(event_id);
CREATE INDEX idx_event_organizers_user_id ON event_organizers(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event_organizers;
DROP TYPE IF EXISTS organizer_role;
-- +goose StatementEnd