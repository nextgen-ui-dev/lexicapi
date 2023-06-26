package article

import (
	"errors"
	"strings"

	"github.com/oklog/ulid/v2"
)

var (
	ErrArticleCategoryNameTooLong = errors.New("Article category can't be longer than 100 characters")
	ErrArticleCategoryNameEmpty   = errors.New("Article category can't be empty")
	ErrInvalidArticleCategoryId   = errors.New("Invalid article category id")
)

func validateArticleCategoryId(idStr string) (id ulid.ULID, err error) {
	id, err = ulid.Parse(idStr)
	if err != nil {
		return id, ErrInvalidArticleCategoryId
	}

	return id, nil
}

func validateArticleCategoryName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrArticleCategoryNameEmpty
	}
	if len(name) > 100 {
		return ErrArticleCategoryNameTooLong
	}

	return nil
}
