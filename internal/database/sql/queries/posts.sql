-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, NOW(), NOW(), $2, $3, $4, $5, $6) RETURNING *;

-- name: GetPostsForUser :many
WITH followed_feeds AS (
  SELECT feed_follows.feed_id
  FROM feed_follows
  WHERE feed_follows.user_id = $1
),
feed_details AS (
  SELECT posts.*, feeds.name AS feed_name
  FROM posts
  JOIN feeds ON posts.feed_id = feeds.id
  WHERE posts.feed_id IN (SELECT feed_id FROM followed_feeds)
)
SELECT * FROM feed_details
ORDER BY published_at 
LIMIT $2;
