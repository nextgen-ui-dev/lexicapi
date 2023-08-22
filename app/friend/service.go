package friend

import "context"

func sendFriendRequest(ctx context.Context, requesterIdStr, requesteeIdStr string) (friend Friend, errs map[string]error, err error) {
	friend, errs, err = NewFriendRequest(requesterIdStr, requesteeIdStr)
	if errs != nil {
		return
	}
	if err != nil {
		return
	}

	return friend, nil, nil
}
