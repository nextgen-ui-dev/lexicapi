package assistant

import "context"

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
