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

var (
	ChatbotInputPrompt = `I'm a helpful chatbot! Ask me anything:`

	ChatbotPiratePrompt         = `Ahoy! Got questions? Spit 'em out, ye scallywag!`
	ChatbotSystemPirateTemplate = `Respond to the following like a pirate: %s`

	ChatbotSpanishPrompt                         = `Enter a phrase, and I'll translate it into Spanish!`
	ChatbotSystemSpanishTranslateTemplate        = `Translate the following phrase into Spanish: %s`
	ChatbotSystemSpanishTranslateFocusedTemplate = `Translate the following phrase into Spanish: %s. 
	Answer with just the translation and no other introductory text.`
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

	fmt.Fprintln(W, ChatbotSpanishPrompt)

	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(W, "reading standard input %v", err)
		os.Exit(1)
	}

	client := openai.NewClient(
		option.WithAPIKey(apikey),
	)

	spanishTranslator := newSpanishTranslator(ctx,
		client,
		*temp,
		MaxOutputTokens,
		openai.ChatModelGPT5_4Mini2026_03_17,
	)

	rawInputPrompt := scanner.Text()
	augmentedPrompt := fmt.Sprintf(ChatbotSystemSpanishTranslateFocusedTemplate, rawInputPrompt)

	response, err := spanishTranslator(augmentedPrompt)
	if err != nil {
		fmt.Fprintf(W, "failed to translate: %v", err)
		os.Exit(1)
	}

	fmt.Fprintf(W, "Response from model (with temperature of %0.2f):\n%s", *temp, response)
}
