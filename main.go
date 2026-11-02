package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"github.com/spf13/pflag"
)

const (
	MaxOutputTokens int64 = 250
)

func main() {
	temp := pflag.Float64("temp", 0.5, "set the model's temperature")
	pflag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	apikey := os.Getenv("OPENAI_API_KEY")
	if apikey == "" {
		log.Println("openai api key is not set")
		os.Exit(1)
	}

	ctx := context.Background()
	client := openai.NewClient(
		option.WithAPIKey(apikey),
	)

	prompt := "Who was the first person to land on the Moon?"

	resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Input:           responses.ResponseNewParamsInputUnion{OfString: openai.String(prompt)},
		Model:           openai.ChatModelGPT5_4Mini2026_03_17,
		Temperature:     openai.Float(*temp),
		MaxOutputTokens: openai.Int(MaxOutputTokens),
	})
	if err != nil {
		log.Printf("error getting response from model: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Response from model (with temperature of %0.2f):\n%s", *temp, resp.OutputText())
}
