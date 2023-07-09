package assistant

import (
	"strings"

	"github.com/jellydator/validation"
)

var (
	ErrSimplificationOriginalTextEmpty   = validation.NewError("assistant:original_text_empty", "Original text can't be empty")
	ErrSimplificationOriginalTextTooLong = validation.NewError("assistant:original_text_too_long", "Original text can't be longer than 8000 characters")
	ErrSimplificationSimplifiedTextEmpty = validation.NewError("assistant:simplified_text_empty", "Simplified text can't be empty")
)

func validateSimplicationOriginalText(originalText string) (err error) {
	originalText = strings.TrimSpace(originalText)
	return validation.Validate(
		&originalText,
		validation.Required.ErrorObject(ErrSimplificationOriginalTextEmpty),
		validation.Length(1, 8000).ErrorObject(ErrSimplificationOriginalTextTooLong),
	)
}

func validateSimplificationSimplifiedText(simplifiedText string) (err error) {
	simplifiedText = strings.TrimSpace(simplifiedText)
	return validation.Validate(
		&simplifiedText,
		validation.Required.ErrorObject(ErrSimplificationSimplifiedTextEmpty),
	)
}
