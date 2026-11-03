package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/spf13/pflag"
)

var (
	DefaultQuery    string = "Who was the first person to land on the Moon?"
	MaxOutputTokens int64  = 250
	W                      = os.Stdout
)

var (
	ChatbotInputPrompt          = `I'm a helpful chatbot! Ask me anything:`
	ChatbotPiratePrompt         = `Ahoy! Got questions? Spit 'em out, ye scallywag!`
	ChatbotSystemPirateTemplate = `Respond to the following like a pirate: %s`
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

	fmt.Fprintln(W, ChatbotPiratePrompt)

	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(W, "reading standard input %v", err)
		os.Exit(1)
	}

	rawInputPrompt := scanner.Text()
	augmentedPrompt := fmt.Sprintf(ChatbotSystemPirateTemplate, rawInputPrompt)

	ctx := context.Background()
	client := openai.NewClient(
		option.WithAPIKey(apikey),
	)

	resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Input:           responses.ResponseNewParamsInputUnion{OfString: openai.String(augmentedPrompt)},
		Model:           openai.ChatModelGPT5_4Mini2026_03_17,
		Temperature:     openai.Float(*temp),
		MaxOutputTokens: openai.Int(MaxOutputTokens),
	})
	if err != nil {
		fmt.Fprintf(W, "error getting response from model: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(W, "Response from model (with temperature of %0.2f):\n%s", *temp, resp.OutputText())
}
