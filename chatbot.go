package main

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
)

type chatClient struct {
	llmClient            *openai.Client
	modelChoice          shared.ChatModel
	modelTemperature     float64
	modelMaxOutputTokens int64
}

var (
	SystemMessage    string = "You are an AI assistant who always talks like a pirate."
	AssistantMessage string = "Arrgh, how can I help you, matey?"
)

func newUserMessage(query string) responses.ResponseInputItemUnionParam {
	return responses.ResponseInputItemParamOfMessage(query, responses.EasyInputMessageRoleUser)
}

func newAssistantMessage(outputFromLLM string) responses.ResponseInputItemUnionParam {
	return responses.ResponseInputItemParamOfMessage(outputFromLLM, responses.EasyInputMessageRoleAssistant)
}

func newSystemInstruction(instruction string) responses.ResponseInputItemUnionParam {
	return responses.ResponseInputItemParamOfMessage(
		instruction,
		responses.EasyInputMessageRoleDeveloper,
	)
}

func (c *chatClient) sendMessage(ctx context.Context, prompt string) (*responses.Response, error) {
	resp, err := c.llmClient.Responses.New(ctx, responses.ResponseNewParams{
		Input:           responses.ResponseNewParamsInputUnion{OfString: openai.String(prompt)},
		Model:           c.modelChoice,
		Temperature:     openai.Float(c.modelTemperature),
		MaxOutputTokens: openai.Int(c.modelMaxOutputTokens),
		Store:           openai.Bool(false),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting response from model: %w", err)
	}

	return resp, nil
}

func (c *chatClient) sendMessageWithHistory(ctx context.Context, history responses.ResponseInputParam) (*responses.Response, error) {
	resp, err := c.llmClient.Responses.New(ctx, responses.ResponseNewParams{
		Model:           c.modelChoice,
		Temperature:     openai.Float(c.modelTemperature),
		MaxOutputTokens: openai.Int(c.modelMaxOutputTokens),
		Instructions:    openai.String(SystemMessage),
		Store:           openai.Bool(false),
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: history,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error getting response with history from model: %w", err)
	}

	return resp, nil
}

func outputAsInput(output []responses.ResponseOutputItemUnion) []responses.ResponseInputItemUnionParam {
	input := make([]responses.ResponseInputItemUnionParam, 0, len(output))
	for _, item := range output {
		var converted responses.ResponseInputItemUnion
		if err := json.Unmarshal([]byte(item.RawJSON()), &converted); err != nil {
			panic(err)
		}
		input = append(input, converted.ToParam())
	}
	return input
}
