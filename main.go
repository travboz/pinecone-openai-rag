package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/spf13/pflag"
)

var (
	TextColourUser      = "\033[34m"
	TextColourAssistant = "\033[32m"
)

var (
	MaxOutputTokens int64     = 250
	W               io.Writer = os.Stdout
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	temp := pflag.Float64("temp", 0.5, "set the model's temperature")
	pflag.Parse()

	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(W, err)
		os.Exit(1)
	}

	apikey := os.Getenv("OPENAI_API_KEY")
	if apikey == "" {
		fmt.Fprintln(W, "openai api key is not set")
		os.Exit(1)
	}

	// Chatbot creation:
	client := openai.NewClient(
		option.WithAPIKey(apikey),
	)

	ctx := context.TODO()

	chatbot := chatClient{
		llmClient:            &client,
		modelChoice:          openai.ChatModelGPT5_4Mini2026_03_17,
		modelTemperature:     *temp,
		modelMaxOutputTokens: MaxOutputTokens,
	}

	scanner := bufio.NewScanner(os.Stdin)

	initialAssistantMessage := TextColourAssistant + "Assistant: How can I help you today?\n"
	fmt.Print(initialAssistantMessage)

	fmt.Print(TextColourUser + "User: ")
	scanner.Scan()
	userInput := scanner.Text()

	// conversationHistory += ChatbotHelpfulAssistantInitial + userInput
	conversationHistory := []responses.ResponseInputItemUnionParam{
		newSystemInstruction(SystemMessage),
		newAssistantMessage(initialAssistantMessage),
		newUserMessage(userInput),
	}

	// Chat loop:
	for {
		if userInput != "exit" {
			break
		}

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
