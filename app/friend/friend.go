package friend

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrCantSendFriendRequestToSelf = errors.New("can't send a friend request to yourself")
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
