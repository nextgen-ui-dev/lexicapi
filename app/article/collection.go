package article

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type CollectionVisibility string

const (
	PRIVATE CollectionVisibility = "private"
	SHARED  CollectionVisibility = "shared"
	PUBLIC  CollectionVisibility = "public"
)

type Collection struct {
	Id         ulid.ULID            `json:"id"`
	CreatorId  ulid.ULID            `json:"creator_id"`
	Name       string               `json:"name"`
	Visibility CollectionVisibility `json:"visibility"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  null.Time            `json:"updated_at"`
	DeletedAt  null.Time            `json:"deleted_at"`
}

func NewCollection(creatorIdStr, name, visibilityStr string) (Collection, map[string]error) {
	errs := make(map[string]error)

	creatorId, err := validateCollectionCreatorId(creatorIdStr)
	if err != nil {
		errs["creator_id"] = err
	}

	if err := validateCollectionName(name); err != nil {
		errs["name"] = err
	}

	visibility, err := validateCollectionVisibility(visibilityStr)
	if err != nil {
		errs["visibility"] = err
	}

	if len(errs) != 0 {
		return Collection{}, errs
	}

	return Collection{
		Id: ulid.Make(),
		CreatorId: creatorId,
		Name: name,
		Visibility: visibility,
		CreatedAt: time.Now(),
	}, nil
}
