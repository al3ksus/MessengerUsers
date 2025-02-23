package friendrequestspsql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/al3ksus/messengerusers/internal/domain/models"
	"github.com/al3ksus/messengerusers/internal/repositories/psql"
	userspsql "github.com/al3ksus/messengerusers/internal/repositories/psql/users"
	"github.com/lib/pq"
)

type FriendRequestsRepository struct {
	db *sql.DB
}

var (
	ErrAlreadyFriends        = errors.New("users are already friends")
	ErrFriendRequestNotFound = errors.New("friend request not found")
)

func NewFriendRequestsRepository(db *sql.DB) *FriendRequestsRepository {
	return &FriendRequestsRepository{
		db: db,
	}
}

// SaveFriendRequest сохраняет запрос на добавление в друзья в базу данных.
func (r *FriendRequestsRepository) Save(ctx context.Context, senderId, recipientId int64) (int64, error) {
	const op = "psql.SaveFriendRequest"

	var id int64
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO friend_requests (
			sender_id, recipient_id, accepted
		) VALUES ($1, $2, false) RETURNING id`,
		senderId, recipientId)

	if err := row.Scan(&id); err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Name() == psql.CodeConstraintUnique {
			return 0, fmt.Errorf("%s, %w", op, ErrAlreadyFriends)
		}

		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Name() == psql.CodeConstraintForeignKey {
			return 0, fmt.Errorf("%s, %w", op, userspsql.ErrUserNotFound)
		}

		return 0, fmt.Errorf("%s, %w", op, err)
	}

	return id, nil
}

func (r *FriendRequestsRepository) Accept(ctx context.Context, requestId int64) error {
	const op = "psql.AcceptFriendRequest"

	res, err := r.db.ExecContext(ctx,
		`UPDATE friend_requests SET accepted = true WHERE id = $1`,
		requestId)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}

	if res, err := res.RowsAffected(); err != nil || res == 0 {
		return fmt.Errorf("%s, %w", op, ErrFriendRequestNotFound)
	}

	return nil
}

func (r *FriendRequestsRepository) Delete(ctx context.Context, requestId int64) error {
	const op = "psql.DeleteFriendRequest"

	res, err := r.db.ExecContext(ctx,
		`DELETE FROM friend_requests WHERE id = $1`,
		requestId)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}

	if count, err := res.RowsAffected(); err != nil || count == 0 {
		return fmt.Errorf("%s, %w", op, ErrFriendRequestNotFound)
	}

	return nil
}

func (r *FriendRequestsRepository) GetFriends(ctx context.Context, userId int64) ([]models.Friend, error) {
	const op = "psql.GetFriends"

	rows, err := r.db.QueryContext(ctx,
		`WITH friends AS (
			SELECT id,
			CASE
  				WHEN sender_id = $1 THEN recipient_id
  				WHEN recipient_id = $1 THEN sender_id
			END AS friend_id,
			accepted
			FROM friend_requests
		)

		SELECT f.id, f.accepted, u.username FROM friends f 
		JOIN users u ON u.id = f.friend_id 
		WHERE friend_id IS NOT NULL;`,
		userId)

	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	defer rows.Close()

	var friends []models.Friend
	for rows.Next() {
		var friend models.Friend
		if err := rows.Scan(&friend.FriendRequestId, &friend.Accepted, &friend.FriendName); err != nil {
			return nil, fmt.Errorf("%s, %w", op, err)
		}

		friends = append(friends, friend)
	}

	return friends, nil
}
