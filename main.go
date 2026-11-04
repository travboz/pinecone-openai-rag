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
	"github.com/spf13/pflag"
)

var (
	DefaultQuery    string          = "Who was the first person to land on the Moon?"
	MaxOutputTokens int64           = 250
	W               io.Writer       = os.Stdout
	ctx             context.Context = context.Background()

	conversationHistory string = ""
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
	chatbot := newInteractiveChatbot(
		ctx,
		client,
		*temp,
		MaxOutputTokens,
		openai.ChatModelGPT5_4Mini2026_03_17,
	)

	scanner := bufio.NewScanner(os.Stdin)

	assistantMessage := "Assistant: How can I help you today?\n\nUser: "
	fmt.Print(assistantMessage)
	scanner.Scan()
	userInput := scanner.Text()

	conversationHistory += ChatbotHelpfulAssistantInitial + userInput

	// Chat loop:
	for userInput != "exit" {
		// Chat interaction:
		response, err := chatbot(conversationHistory)
		if err != nil {
			fmt.Fprintf(W, "generation failed: %v", err)
			log.Println(conversationHistory)
			os.Exit(1)
		}

		llmResponseText := fmt.Sprintf("\nAssistant: %s", response)
		fmt.Fprintln(W, llmResponseText)

		fmt.Fprint(W, "\nUser: ")
		scanner.Scan()
		userInput = scanner.Text()

		// Append prompt and response to history.
		conversationHistory += fmt.Sprintf("%s\nUser: %s", llmResponseText, userInput)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(W, "scanner error: %v", err)
		os.Exit(1)
	}
}
