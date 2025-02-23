package friendrequests

import (
	"context"
	"errors"
	"fmt"

	"github.com/al3ksus/messengerusers/internal/domain/models"
	"github.com/al3ksus/messengerusers/internal/logger"
	friendrequestspsql "github.com/al3ksus/messengerusers/internal/repositories/psql/friendrequests"
	userspsql "github.com/al3ksus/messengerusers/internal/repositories/psql/users"
)

type FriendRequests struct {
	log                   logger.Logger
	friendRequestSaver    FriendRequestSaver
	friendRequestChanger  FriendRequestChanger
	friendRequestProvider FriendRequestProvider
}

type FriendRequestSaver interface {
	Save(ctx context.Context, senderId, recipientId int64) (int64, error)
}

type FriendRequestProvider interface {
	GetFriends(ctx context.Context, userId int64) ([]models.Friend, error)
}

type FriendRequestChanger interface {
	Accept(ctx context.Context, requestId int64) error
	Delete(ctx context.Context, requestId int64) error
}

var (
	ErrAlreadyFriends        = errors.New("users are already friends")
	ErrFriendRequestNotFound = errors.New("friend request not found")
	ErrUserNotFound          = errors.New("user not found")
)

func NewFriendRequests(
	log logger.Logger,
	friendRequestSaver FriendRequestSaver,
	friendRequestChanger FriendRequestChanger,
	friendRequestProvider FriendRequestProvider,
) *FriendRequests {
	return &FriendRequests{
		log:                   log,
		friendRequestSaver:    friendRequestSaver,
		friendRequestChanger:  friendRequestChanger,
		friendRequestProvider: friendRequestProvider,
	}
}

func (f *FriendRequests) Send(ctx context.Context, sender_id, recipient_id int64) (int64, error) {
	const op = "friendrequests.Send"

	id, err := f.friendRequestSaver.Save(ctx, sender_id, recipient_id)
	if err != nil {
		if errors.Is(err, friendrequestspsql.ErrAlreadyFriends) {
			f.log.Warnf("%s: %v", op, err)
			return 0, fmt.Errorf("%s, %w", op, ErrAlreadyFriends)
		}

		if errors.Is(err, userspsql.ErrUserNotFound) {
			f.log.Warnf("%s: %v", op, err)
			return 0, fmt.Errorf("%s, %w", op, ErrUserNotFound)
		}

		f.log.Errorf("%s: error sending friend request senderId=%d, recipient_id=%d, %v", op, sender_id, recipient_id, err)
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	return id, nil
}

func (f *FriendRequests) Accept(ctx context.Context, requestId int64) error {
	const op = "friendrequests.Accept"

	err := f.friendRequestChanger.Accept(ctx, requestId)
	if err != nil {
		if errors.Is(err, friendrequestspsql.ErrFriendRequestNotFound) {
			f.log.Warnf("%s: %v", op, err)
			return fmt.Errorf("%s, %w", op, ErrFriendRequestNotFound)
		}

		f.log.Errorf("%s: error accepting friend request, requestId=%d, %v", op, requestId, err)
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}

func (f *FriendRequests) Delete(ctx context.Context, requestId int64) error {
	const op = "friendrequests.Delete"

	err := f.friendRequestChanger.Delete(ctx, requestId)
	if err != nil {
		if errors.Is(err, friendrequestspsql.ErrFriendRequestNotFound) {
			f.log.Warnf("%s: %v", op, err)
			return fmt.Errorf("%s, %w", op, ErrFriendRequestNotFound)
		}

		f.log.Errorf("%s: error deleting friend request, requestId=%d, %v", op, requestId, err)
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}

func (f *FriendRequests) GetFriends(ctx context.Context, userId int64) ([]models.Friend, error) {
	const op = "friendrequests.GetFriends"

	friends, err := f.friendRequestProvider.GetFriends(ctx, userId)
	if err != nil {
		f.log.Errorf("%s: error getting friends, userId=%d, %v", op, userId, err)
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return friends, nil
}
