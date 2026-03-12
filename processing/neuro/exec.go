package neuro

import (
	"TgBotUltimate/server/routes/helper"
	"TgBotUltimate/types/Neuro"
	"context"
	"fmt"
	"os"
)

func Ask(ctx context.Context, prompt string) string {
	return ""
}

func Parameters(ctx context.Context, prompt string) (*Neuro.Response, error) {
	var response Neuro.Response
	err := helper.Post(
		ctx,
		fmt.Sprintf("http://localhost:%s/parse", os.Getenv("NEURO_PORT")),
		nil,
		Neuro.Request{Text: prompt},
		&response,
	)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
