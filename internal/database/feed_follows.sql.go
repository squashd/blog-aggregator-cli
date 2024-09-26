// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: feed_follows.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO
        feed_follows (id, created_at, updated_at, feed_id, user_id)
    VALUES
        ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at, feed_id, user_id
)
SELECT
    inserted_feed_follow.id, inserted_feed_follow.created_at, inserted_feed_follow.updated_at, inserted_feed_follow.feed_id, inserted_feed_follow.user_id,
    feeds.name AS feed_name,
    users.name AS user_name
FROM
    inserted_feed_follow
    INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
    INNER JOIN users ON inserted_feed_follow.user_id = users.id
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	FeedID    uuid.UUID
	UserID    uuid.UUID
}

type CreateFeedFollowRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	FeedID    uuid.UUID
	UserID    uuid.UUID
	FeedName  string
	UserName  string
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (CreateFeedFollowRow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.FeedID,
		arg.UserID,
	)
	var i CreateFeedFollowRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FeedID,
		&i.UserID,
		&i.FeedName,
		&i.UserName,
	)
	return i, err
}

const deleteFeedFollow = `-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE feed_id = $1 AND user_id = $2
`

type DeleteFeedFollowParams struct {
	FeedID uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteFeedFollow(ctx context.Context, arg DeleteFeedFollowParams) error {
	_, err := q.db.ExecContext(ctx, deleteFeedFollow, arg.FeedID, arg.UserID)
	return err
}

const getFollowsForUser = `-- name: GetFollowsForUser :many
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
    INNER JOIN feed_follows ON users_data.id = feed_follows.user_id
`

type GetFollowsForUserRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	FeedID    uuid.UUID
	UserID    uuid.UUID
	FeedName  string
	Name      string
}

func (q *Queries) GetFollowsForUser(ctx context.Context, name string) ([]GetFollowsForUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFollowsForUser, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetFollowsForUserRow{}
	for rows.Next() {
		var i GetFollowsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedID,
			&i.UserID,
			&i.FeedName,
			&i.Name,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
