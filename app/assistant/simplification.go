package assistant

type Simplication struct {
	OriginalText   string `json:"original_text"`
	SimplifiedText string `json:"simplified_text"`
}

func NewSimplification(originalText string) (s Simplication, err error) {
	if err = validateSimplicationOriginalText(originalText); err != nil {
		return
	}

	return Simplication{OriginalText: originalText}, nil
}

func (s *Simplication) Simplify(simplifiedText string) (err error) {
	if err = validateSimplificationSimplifiedText(simplifiedText); err != nil {
		return
	}

	s.SimplifiedText = simplifiedText
	return nil
}
