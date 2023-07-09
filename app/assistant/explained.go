package assistant

type Explained struct {
	Text        string `json:"text"`
	Explanation string `json:"explanation"`
}

func NewExplained(text string) (explained Explained, err error) {
	if err = validateExplainedText(text); err != nil {
		return
	}

	return Explained{Text: text}, nil
}

func (e *Explained) Explain(explanation string) (err error) {
	if err = validateExplainedExplanation(explanation); err != nil {
		return
	}

	e.Explanation = explanation

	return nil
}
