-- +goose Up
CREATE TABLE refresh_tokens(
	token TEXT PRIMARY KEY NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	expires_at TIMESTAMP NOT NULL,
	revoked_at TIMESTAMP DEFAULT NULL,

	user_id UUID NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE refresh_tokens;
