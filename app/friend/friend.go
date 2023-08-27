package friend

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrCantSendFriendRequestToSelf             = errors.New("can't send a friend request to yourself")
	ErrFriendRequestAlreadyRejected            = errors.New("friend request already rejected")
	ErrFriendRequestAlreadyAccepted            = errors.New("friend request already accepted")
	ErrCantAcceptFriendRequestOfOtherRequestee = errors.New("can't accept a friend request that belongs to other requestee")
	ErrCantRejectFriendRequestOfOtherRequestee = errors.New("can't reject a friend request that belongs to other requestee")
	ErrCantUnfriendOtherUserFriend             = errors.New("can't unfriend other user's friend")
)

type FriendStatus string

const (
	PENDING  FriendStatus = "pending"
	FRIENDED FriendStatus = "friended"
	REJECTED FriendStatus = "rejected"
)

type Friend struct {
	Id          ulid.ULID    `json:"id"`
	RequesterId ulid.ULID    `json:"requester_id"`
	RequesteeId ulid.ULID    `json:"requestee_id"`
	Status      FriendStatus `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	DeletedAt   null.Time    `json:"deleted_at"`
}

func NewFriendRequest(requesterIdStr, requesteeIdStr string) (Friend, map[string]error, error) {
	errs := make(map[string]error)

	requesterId, err := validateRequesterId(requesterIdStr)
	if err != nil {
		errs["requester_id"] = err
	}

	requesteeId, err := validateRequesteeId(requesteeIdStr)
	if err != nil {
		errs["requestee_id"] = err
	}

	// requesterId == requesteeId
	if requesterId.Compare(requesteeId) == 0 {
		return Friend{}, nil, ErrCantSendFriendRequestToSelf
	}

	if len(errs) != 0 {
		return Friend{}, errs, nil
	}

	return Friend{
		Id:          ulid.Make(),
		RequesterId: requesterId,
		RequesteeId: requesteeId,
		Status:      PENDING,
		CreatedAt:   time.Now(),
	}, nil, nil
}

func (f *Friend) AcceptFriendRequest(requesteeId ulid.ULID) (err error) {
	if f.Status == REJECTED || f.DeletedAt.Valid {
		return ErrFriendRequestAlreadyRejected
	} else if f.Status == FRIENDED {
		return ErrFriendRequestAlreadyAccepted
	}

	if f.RequesteeId.Compare(requesteeId) != 0 {
		return ErrCantAcceptFriendRequestOfOtherRequestee
	}

	f.Status = FRIENDED

	return nil
}

func (f *Friend) RejectFriendRequest(requesteeId ulid.ULID) (err error) {
	if f.Status == REJECTED || f.DeletedAt.Valid {
		return ErrFriendRequestAlreadyRejected
	} else if f.Status == FRIENDED {
		return ErrFriendRequestAlreadyAccepted
	}

	if f.RequesteeId.Compare(requesteeId) != 0 {
		return ErrCantRejectFriendRequestOfOtherRequestee
	}

	f.Status = REJECTED
	f.DeletedAt = null.TimeFrom(time.Now())

	return nil
}

func (f *Friend) Unfriend(userId ulid.ULID) (err error) {
	if f.Status == REJECTED && f.DeletedAt.Valid {
		return ErrFriendRequestAlreadyRejected
	} else if f.Status == FRIENDED {
		return ErrFriendRequestAlreadyAccepted
	}

	if f.RequesteeId.Compare(userId) != 0 && f.RequesterId.Compare(userId) != 0 {
		return ErrCantUnfriendOtherUserFriend
	}

	f.DeletedAt = null.TimeFrom(time.Now())

	return nil
}
