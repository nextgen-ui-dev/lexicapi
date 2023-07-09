package assistant

import "context"

func simplifyText(ctx context.Context, originalText string) (s Simplication, err error) {
	s, err = NewSimplification(originalText)
	if err != nil {
		return
	}

	simplifiedText, err := generateSimplifiedText(ctx, originalText)
	if err != nil {
		return
	}

	if err = s.Simplify(simplifiedText); err != nil {
		return
	}

	return s, nil
}

func explainText(ctx context.Context, body ExplainTextReq) (explained Explained, err error) {
	explained, err = NewExplained(body.Text)
	if err != nil {
		return
	}

	explanation, err := generateTextExplanation(ctx, body.Text)
	if err != nil {
		return
	}

	if err = explained.Explain(explanation); err != nil {
		return
	}

	return explained, nil
}
