package friend

import (
	"github.com/jellydator/validation"
	"github.com/oklog/ulid/v2"
)

var (
	ErrInvalidRequesterId = validation.NewError("friend:invalid_requester_id", "Invalid requester id")
	ErrInvalidRequesteeId = validation.NewError("friend:invalid_requestee_id", "Invalid requestee id")
)

func validateRequesterId(idStr string) (id ulid.ULID, err error) {
	id, err = ulid.Parse(idStr)
	if err != nil {
		return id, ErrInvalidRequesterId
	}

	return id, nil
}

func validateRequesteeId(idStr string) (id ulid.ULID, err error) {
	id, err = ulid.Parse(idStr)
	if err != nil {
		return id, ErrInvalidRequesterId
	}

	return id, nil
}
