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

func acceptFriendRequest(ctx context.Context, requesteeIdStr, friendIdStr string) (friend Friend, err error) {
	requesteeId, err := validateRequesteeId(requesteeIdStr)
	if err != nil {
		return
	}

	friendId, err := validateFriendId(friendIdStr)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to accept friend request")
		return
	}

	defer tx.Rollback(ctx)

	friend, err = findFriendById(ctx, tx, friendId)
	if err != nil {
		return
	}

	if err = friend.AcceptFriendRequest(requesteeId); err != nil {
		return
	}

	friend, err = updateFriend(ctx, tx, friend)
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to accept friend request")
		return
	}

	return friend, nil
}

func rejectFriendRequest(ctx context.Context, requesteeIdStr, friendIdStr string) (friend Friend, err error) {
	requesteeId, err := validateRequesteeId(requesteeIdStr)
	if err != nil {
		return
	}

	friendId, err := validateFriendId(friendIdStr)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to reject friend request")
		return
	}

	defer tx.Rollback(ctx)

	friend, err = findFriendById(ctx, tx, friendId)
	if err != nil {
		return
	}

	if err = friend.RejectFriendRequest(requesteeId); err != nil {
		return
	}

	friend, err = updateFriend(ctx, tx, friend)
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to reject friend request")
		return
	}

	return friend, nil
}

func Unfriend(ctx context.Context, userIdStr, friendIdStr string) (friend Friend, err error) {
	userId, err := validateRequesteeId(userIdStr)
	if err != nil {
		return
	}

	friendId, err := validateFriendId(friendIdStr)
	if err != nil {
		return
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to unfriend")
		return
	}

	defer tx.Rollback(ctx)

	friend, err = findFriendById(ctx, tx, friendId)
	if err != nil {
		return
	}

	if err = friend.Unfriend(userId); err != nil {
		return
	}

	friend, err = updateFriend(ctx, tx, friend)
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		log.Err(err).Msg("Failed to unfriend")
		return
	}

	return friend, nil
}
