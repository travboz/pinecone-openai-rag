package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
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
)

func main() {
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

	fmt.Fprint(W, ChatbotHelpfulAssistantInitial)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		if input == "exit" {
			break
		}

		// Chat interaction:
		rawInputPrompt := scanner.Text()
		// augmentedPrompt := fmt.Sprintf(ChatbotSystemSpanishTranslateFocusedTemplate, rawInputPrompt)

		response, err := chatbot(rawInputPrompt)
		if err != nil {
			fmt.Fprintf(W, "failed to translate: %v", err)
			os.Exit(1)
		}

		fmt.Fprintf(W, "\nAssistant (temp=%0.1f): %s\n", *temp, response)
		fmt.Fprint(W, "\nUser: ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(W, "reading standard input %v", err)
		os.Exit(1)
	}

}
