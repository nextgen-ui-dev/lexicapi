package article

import (
	"errors"
	"strings"
)

var (
	ErrArticleCategoryNameTooLong = errors.New("Article category can't be longer than 100 characters")
	ErrArticleCategoryNameEmpty   = errors.New("Article category can't be empty")
)

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
