package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

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
	// TODO: Fix this, needs new system message as cannot use this var.
	SystemMessage string = "You are an AI customer support technician who is knowledgeable about software products created by the company called GROSS. One such product is a web browser called Flamehamster. You are to answer user queries below solely on the following documentation: %s"

// AssistantMessage string = "How can I help you today?"
)

func newUserMessage(query string) responses.ResponseInputItemUnionParam {
	return responses.ResponseInputItemParamOfMessage(query, responses.EasyInputMessageRoleUser)
}

func newAssistantMessage(outputFromLLM string) responses.ResponseInputItemUnionParam {
	return responses.ResponseInputItemParamOfMessage(outputFromLLM, responses.EasyInputMessageRoleAssistant)
}

func newSystemInstructionWithKnowledge(instruction, knowledge string) responses.ResponseInputItemUnionParam {
	augmentedInstruction := fmt.Sprintf(instruction, knowledge)

	return responses.ResponseInputItemParamOfMessage(
		augmentedInstruction,
		responses.EasyInputMessageRoleDeveloper,
	)
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
		// Instructions:    openai.String(SystemMessage), // TODO: Fix this, needs new system message as cannot use var declared at top of file.
		Store: openai.Bool(false),
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

// TODO: Make this into an actual, workable function. Currently: PLACEHOLDER for reference.
func HandleChatInteraction(ctx context.Context, chatbot chatClient, flamehamsterFile []byte) {
	scanner := bufio.NewScanner(os.Stdin)

	initialAssistantMessage := TextColourAssistant + "Assistant: How can I help you today?\n"
	fmt.Print(initialAssistantMessage)

	fmt.Print(TextColourUser + "User: ")
	scanner.Scan()
	userInput := scanner.Text()

	// conversationHistory += ChatbotHelpfulAssistantInitial + userInput
	conversationHistory := []responses.ResponseInputItemUnionParam{
		newSystemInstructionWithKnowledge(SystemMessage, string(flamehamsterFile)),
		newAssistantMessage(initialAssistantMessage),
		newUserMessage(userInput),
	}

	// Chat loop:
	for userInput != "exit" {

		// Chat interaction:
		response, err := chatbot.sendMessageWithHistory(ctx, conversationHistory)
		if err != nil {
			fmt.Fprintf(W, "generation failed: %v", err)
			log.Println(conversationHistory)
			os.Exit(1)
		}

		llmResponseText := fmt.Sprintf(TextColourAssistant+"Assistant: %s\n", response.OutputText())
		fmt.Fprint(W, llmResponseText)

		fmt.Fprint(W, TextColourUser+"User: ")
		scanner.Scan()
		userInput = scanner.Text()

		// Append LLM output response to history.
		conversationHistory = append(
			conversationHistory,
			outputAsInput(response.Output)...,
		)

		// Append user input prompt to history.
		conversationHistory = append(
			conversationHistory,
			newUserMessage(userInput),
		)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(W, "scanner error: %v", err)
		os.Exit(1)
	}
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
