package article

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gopkg.in/guregu/null.v4"
)

type CollectionVisibility string

const (
	PRIVATE CollectionVisibility = "private"
	SHARED CollectionVisibility = "shared"
	PUBLIC CollectionVisibility = "public"
)

type Collection struct {
	Id        ulid.ULID `json:"id"`
	CreatorId ulid.ULID `json:"creator_id"`
	Name string `json:"name"`
	Visibility CollectionVisibility `json:"visibility"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    null.Time   `json:"updated_at"`
	DeletedAt    null.Time   `json:"deleted_at"`
}