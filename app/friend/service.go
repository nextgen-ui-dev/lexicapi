package friend

import (
	"context"

	"github.com/rs/zerolog/log"
)

func sendFriendRequest(ctx context.Context, requesterIdStr, requesteeIdStr string) (friend Friend, errs map[string]error, err error) {
	friend, errs, err = NewFriendRequest(requesterIdStr, requesteeIdStr)
	if errs != nil {
		return
	}
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to send friend request")
		return
	}

	defer tx.Rollback(ctx)

	existingFriend, err := findPendingOrFriendedFriendByUserIds(ctx, tx, friend.RequesterId, friend.RequesteeId)
	if err == nil {
		if existingFriend.Status == PENDING {
			return Friend{}, nil, ErrFriendRequestAlreadyExist
		}

		return Friend{}, nil, ErrAlreadyFriends
	}

	if err != ErrFriendDoesNotExist {
		return
	}

	friend, err = createFriend(ctx, tx, friend)
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to send friend request")
		return Friend{}, nil, err
	}

	return friend, nil, nil
}
