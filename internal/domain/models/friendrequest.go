package models

// FriendRequest модель данных запроса на добавление в друзья.
type FriendRequest struct {
	Id          int64
	SenderId    int64
	RecipientId int64
	Accepted    bool
}

type Friend struct {
	FriendRequestId int64
	FriendName      string
	Accepted        bool
}
