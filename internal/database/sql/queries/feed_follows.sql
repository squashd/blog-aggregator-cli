-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO
        feed_follows (id, created_at, updated_at, feed_id, user_id)
    VALUES
        ($1, $2, $3, $4, $5) RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM
    inserted_feed_follow
    INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
    INNER JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFollowsForUser :many
WITH users_data AS (
    SELECT
        users.id,
        users.name
    FROM
        users
    WHERE
        users.name = $1
),
feed_follows AS (
    SELECT
        feed_follows.id,
        feed_follows.created_at,
        feed_follows.updated_at,
        feed_follows.feed_id,
        feed_follows.user_id,
        feeds.name AS feed_name
    FROM
        feed_follows
        INNER JOIN feeds ON feed_follows.feed_id = feeds.id
)
SELECT
    feed_follows.id,
    feed_follows.created_at,
    feed_follows.updated_at,
    feed_follows.feed_id,
    feed_follows.user_id,
    feed_follows.feed_name,
    users_data.name
FROM
    users_data
    INNER JOIN feed_follows ON users_data.id = feed_follows.user_id;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE feed_id = $1 AND user_id = $2;
