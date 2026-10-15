-- name: CreatePost :one
INSERT INTO posts (created_at, updated_at, title, url, description, published_at, feed_id)
VALUES(
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7
)
RETURNING *;

-- name: GetPosts :many
SELECT
	posts.*,
	feeds.name AS feed_name
FROM
	posts
	INNER JOIN feeds ON posts.feed_id = feeds.id
	INNER JOIN users ON feeds.user_id = users.id
WHERE
	users.name = $1
ORDER BY posts.published_at ASC NULLS FIRST
LIMIT $2;
