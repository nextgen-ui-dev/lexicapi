package article

import (
	"strings"

	"github.com/jellydator/validation"
	"github.com/oklog/ulid/v2"
)

var (
	ErrInvalidCollectionId         = validation.NewError("article:invalid_collection_id", "invalid collection id")
	ErrInvalidCollectionCreatorId  = validation.NewError("article:invalid_collection_creator_id", "invalid collection creator id")
	ErrEmptyCollectionName         = validation.NewError("article:empty_collection_name", "collection name can't be empty")
	ErrCollectionNameTooLong       = validation.NewError("article:collection_name_too_long", "collection name can't be longer than 100 characters")
	ErrInvalidCollectionVisibility = validation.NewError("article:invalid_collection_visibility", "invalid collection visibility")
)

func validateCollectionId(idStr string) (id ulid.ULID, err error) {
	id, err = ulid.Parse(idStr)
	if err != nil {
		return id, ErrInvalidCollectionId
	}

	return id, nil
}

func validateCollectionCreatorId(idStr string) (id ulid.ULID, err error) {
	id, err = ulid.Parse(idStr)
	if err != nil {
		return id, ErrInvalidCollectionCreatorId
	}

	return id, nil
}

func validateCollectionName(name string) (err error) {
	name = strings.TrimSpace(name)
	return validation.Validate(
		&name,
		validation.Required.ErrorObject(ErrEmptyCollectionName),
		validation.Length(1, 100).ErrorObject(ErrCollectionNameTooLong),
	)
}

func validateCollectionVisibility(visibilityStr string) (visibility CollectionVisibility, err error) {
	switch visibilityStr {
	case string(PRIVATE):
		return PRIVATE, nil
	case string(SHARED):
		return SHARED, nil
	case string(PUBLIC):
		return PUBLIC, nil
	default:
		return CollectionVisibility(""), ErrInvalidCollectionVisibility
	}
}
