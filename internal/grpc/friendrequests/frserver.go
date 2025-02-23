package friendrequestsgrpc

import (
	"context"
	"errors"

	messengerv1 "github.com/al3ksus/messengerprotos/gen/go"
	"github.com/al3ksus/messengerusers/internal/domain/models"
	"github.com/al3ksus/messengerusers/internal/services/friendrequests"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FriendRequestsServerAPI struct {
	messengerv1.UnimplementedFriendRequestsServer
	friendRequests FriendRequests
}

type FriendRequests interface {
	Send(ctx context.Context, sender_id, recipient_id int64) (int64, error)
	Accept(ctx context.Context, requestId int64) error
	Delete(ctx context.Context, requestId int64) error
	GetFriends(ctx context.Context, userId int64) ([]models.Friend, error)
}

func Register(gRPCServer *grpc.Server, friendRequests FriendRequests) {
	messengerv1.RegisterFriendRequestsServer(gRPCServer, &FriendRequestsServerAPI{friendRequests: friendRequests})
}

func (s *FriendRequestsServerAPI) Send(
	ctx context.Context,
	in *messengerv1.SendRequest) (
	*messengerv1.SendResponse, error) {
	id, err := s.friendRequests.Send(ctx, in.GetSenderId(), in.GetRecipientId())
	if err != nil {
		if errors.Is(err, friendrequests.ErrAlreadyFriends) {
			return nil, status.Error(codes.AlreadyExists, "already friends")
		}

		if errors.Is(err, friendrequests.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &messengerv1.SendResponse{FriendRequestId: id}, nil
}

func (s *FriendRequestsServerAPI) Accept(
	ctx context.Context,
	in *messengerv1.AcceptRequest) (
	*messengerv1.AcceptResponse, error) {
	err := s.friendRequests.Accept(ctx, in.GetFriendRequestId())
	if err != nil {
		if errors.Is(err, friendrequests.ErrFriendRequestNotFound) {
			return nil, status.Error(codes.NotFound, "friend request not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &messengerv1.AcceptResponse{}, nil
}

func (s *FriendRequestsServerAPI) Delete(
	ctx context.Context,
	in *messengerv1.DeleteRequest) (
	*messengerv1.DeleteResponse, error) {
	err := s.friendRequests.Delete(ctx, in.GetFriendRequestId())
	if err != nil {
		if errors.Is(err, friendrequests.ErrFriendRequestNotFound) {
			return nil, status.Error(codes.NotFound, "friend request not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &messengerv1.DeleteResponse{}, nil
}

func (s *FriendRequestsServerAPI) GetFriends(
	ctx context.Context,
	in *messengerv1.GetFriendsRequest) (
	*messengerv1.GetFriendsResponse, error) {
	friends, err := s.friendRequests.GetFriends(ctx, in.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	var friendsResponse []*messengerv1.Friend
	for _, friend := range friends {
		friendsResponse = append(friendsResponse, &messengerv1.Friend{
			FriendRequestId: friend.FriendRequestId,
			FriendName:      friend.FriendName,
			Accepted:        friend.Accepted,
		})
	}

	return &messengerv1.GetFriendsResponse{Friend: friendsResponse}, nil
}
