package article

import (
	"errors"
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

var (
	ErrNotAllowedToUpdateCollection = errors.New("not allowed to update collection")
)

type CollectionVisibility struct {
	null.String
}

var (
	PRIVATE CollectionVisibility = CollectionVisibility{null.StringFrom("private")}
	SHARED  CollectionVisibility = CollectionVisibility{null.StringFrom("shared")}
	PUBLIC  CollectionVisibility = CollectionVisibility{null.StringFrom("public")}
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

func (c *Collection) Update(creatorId ulid.ULID, name, visibility null.String) (map[string]error, error) {
	errs := make(map[string]error)

	if c.CreatorId.Compare(creatorId) != 0 {
		return nil, ErrNotAllowedToUpdateCollection
	}

	if name.Valid {
		if err := validateCollectionName(name.String); err != nil {
			errs["name"] = err
		}

		c.Name = name.String
	}

	if visibility.Valid {
		visibility, err := validateCollectionVisibility(visibility.String) 
		if err != nil {
			errs["visibility"] = err
		}

		c.Visibility = visibility
	}

	if len(errs) != 0 {
		return errs, nil
	}

	c.UpdatedAt = null.TimeFrom(time.Now())

	return nil, nil
}
