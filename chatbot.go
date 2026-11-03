package main

import (
	"context"
	"fmt"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
)

func newInteractiveChatbot(
	ctx context.Context,
	llmClient openai.Client,
	temperature float64,
	maxOutputTokens int64,
	model shared.ChatModel,
) func(string) (string, error) {
	return func(prompt string) (string, error) {
		resp, err := llmClient.Responses.New(ctx, responses.ResponseNewParams{
			Input:           responses.ResponseNewParamsInputUnion{OfString: openai.String(prompt)},
			Model:           model,
			Temperature:     openai.Float(temperature),
			MaxOutputTokens: openai.Int(maxOutputTokens),
		})
		if err != nil {
			return "", fmt.Errorf("error getting response from model: %w", err)
		}

		return resp.OutputText(), nil
	}
}
