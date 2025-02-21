-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TYPE channel_overlay_now_playing_preset ADD VALUE 'COMPACT';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
