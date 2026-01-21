package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Follower struct {
	UserID      int64  `json:"user_id"`
	FollowingID int64  `json:"following_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followingID, userID int64) error {
	query := `INSERT INTO followers (user_id, following_id) VALUES ($1, $2)`
	_, err := s.db.ExecContext(ctx, query, userID, followingID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return ErrConflict
			}
		}
		return err
	}
	return nil
}

func (s *FollowerStore) Unfollow(ctx context.Context, followingID, userID int64) error {
	query := `
		DELETE FROM followers
		WHERE user_id = $1 AND following_id = $2
	`

	_, err := s.db.ExecContext(ctx, query, userID, followingID)
	return err
}
