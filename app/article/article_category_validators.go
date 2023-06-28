package article

import (
	"strings"

	"github.com/jellydator/validation"
	"github.com/oklog/ulid/v2"
)

var (
	ErrArticleCategoryNameTooLong = validation.NewError("article:category_name_too_long", "Category name can't be longer than 100 characters")
	ErrArticleCategoryNameEmpty   = validation.NewError("article:category_name_empty", "Article category can't be empty")
	ErrInvalidArticleCategoryId   = validation.NewError("article:invalid_category_id", "Invalid article category id")
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
	return validation.Validate(
		&name,
		validation.Required.ErrorObject(ErrArticleCategoryNameEmpty),
		validation.Length(1, 100).ErrorObject(ErrArticleCategoryNameTooLong),
	)
}
