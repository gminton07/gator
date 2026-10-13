-- name: CreateFeed :one
INSERT INTO feeds (created_at, updated_at, name, url, user_id)
VALUES(
	$1,
	$2,
	$3,
	$4,
	$5
)
RETURNING *;

-- name: GetFeed :one
SELECT
  * 
FROM
  feeds
WHERE
  feeds.url = $1;

-- name: GetFeeds :many
SELECT
  feeds.name,
  feeds.url,
  users.name
FROM
  feeds
  INNER JOIN users on feeds.user_id = users.id;

-- name: MarkFeedFetched :exec
-- ID, TIMESTAMP
UPDATE feeds
SET updated_at = $2, last_fetched_at = $2
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST;
