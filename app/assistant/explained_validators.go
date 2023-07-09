package assistant

import (
	"strings"

	"github.com/jellydator/validation"
)

var (
	ErrExplainedTextEmpty        = validation.NewError("assistant:text_empty", "Text can't be empty")
	ErrExplainedTextTooLong      = validation.NewError("assistant:text_too_long", "Text can't be longer than 255 characters")
	ErrExplainedExplanationEmpty = validation.NewError("assistant:explanation_empty", "Explanation can't be empty")
)

func validateExplainedText(text string) (err error) {
	text = strings.TrimSpace(text)
	return validation.Validate(
		&text,
		validation.Required.ErrorObject(ErrExplainedTextEmpty),
		validation.Length(1, 255).ErrorObject(ErrExplainedTextTooLong),
	)
}

func validateExplainedExplanation(explanation string) (err error) {
	explanation = strings.TrimSpace(explanation)
	return validation.Validate(
		&explanation,
		validation.Required.ErrorObject(ErrExplainedExplanationEmpty),
	)
}
