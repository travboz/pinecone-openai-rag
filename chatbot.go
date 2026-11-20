package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
)

type chatClient struct {
	llmClient            openai.Client
	modelChoice          shared.ChatModel
	modelTemperature     float64
	modelMaxOutputTokens int64
}

func NewChatClient(
	apikey string,
	modelChoice shared.ChatModel,
	temperature float64,
	maxOutputTokens int64,
) *chatClient {

	client := openai.NewClient(
		option.WithAPIKey(apikey),
	)

	// openai.ChatModelGPT5_4Mini2026_03_17
	return &chatClient{
		llmClient:            client,
		modelChoice:          modelChoice,
		modelTemperature:     temperature,
		modelMaxOutputTokens: maxOutputTokens,
	}
}

func augmentUserInputWithRelevantKnowledge(userInput, knowledge, augmentPrompt string) string {
	return fmt.Sprintf(augmentPrompt, knowledge, userInput)
}

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

func (c *chatClient) sendMessageWithHistory(ctx context.Context, history conversationHistory) (*responses.Response, error) {
	resp, err := c.llmClient.Responses.New(ctx, responses.ResponseNewParams{
		Model:           c.modelChoice,
		Temperature:     openai.Float(c.modelTemperature),
		MaxOutputTokens: openai.Int(c.modelMaxOutputTokens),
		Store:           openai.Bool(false),
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: responses.ResponseInputParam(history),
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
		// Reasoning isn't useful across turns here, and with Store:false it comes
		// back without its encrypted payload anyway.
		if item.Type == "reasoning" {
			continue
		}
		raw := item.RawJSON()
		if raw == "" {
			slog.Warn("output item has no raw JSON", slog.String("type", item.Type))
			continue
		}
		var converted responses.ResponseInputItemUnion
		if err := json.Unmarshal([]byte(raw), &converted); err != nil {
			slog.Warn("skipping unconvertible output item", slog.String("type", item.Type), slog.Any("error", err))
			continue
		}
		input = append(input, converted.ToParam())
	}
	return input
}

/* OpenAI chatbot creation reference:
# Retrieve API key for OpenAI API.
Chatbot creation:
client := openai.NewClient(
	option.WithAPIKey(apikey),
)

# Construct a client with basic config.
chatbot := chatClient{
	llmClient:            &client,
	modelChoice:          openai.ChatModelGPT5_4Mini2026_03_17,
	modelTemperature:     temp,
	modelMaxOutputTokens: MaxOutputTokens,
}
*/
