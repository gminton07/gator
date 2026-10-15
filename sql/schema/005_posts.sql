-- +goose Up
CREATE TABLE posts (
	id serial PRIMARY KEY,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	title TEXT NOT NULL,
	url TEXT UNIQUE NOT NULL,
	description TEXT,
	published_at timestamp,
	feed_id serial NOT NULL
		REFERENCES feeds(id)
		ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
