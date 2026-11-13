package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/pinecone-io/go-pinecone/v7/pinecone"
)

var (
	TextColourUser      = "\033[34m"
	TextColourAssistant = "\033[32m"
)

var (
	MaxOutputTokens int64     = 500
	W               io.Writer = os.Stdout

	EmbeddingModel string = "llama-text-embed-v2"
	IndexName      string = "gross-app"
	ManualName     string = "flamehamster"

	temp float64 = 0
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// temp := pflag.Float64("temp", 0.5, "set the model's temperature")
	// pflag.Parse()

	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(W, err)
		os.Exit(1)
	}

	apikey := os.Getenv("OPENAI_API_KEY")
	if apikey == "" {
		fmt.Fprintln(W, "openai api key is not set")
		os.Exit(1)
	}

	pineconeApiKey := os.Getenv("PINECONE_API_KEY")
	if pineconeApiKey == "" {
		fmt.Fprintln(W, "pinecone api key is not set")
		os.Exit(1)
	}

	ctx := context.TODO()

	pc, err := NewPineconeClient(pineconeApiKey)
	if err != nil {
		fmt.Fprintf(W, "failed to create pinecone client: %v\n", err)
		os.Exit(1)
	}

	err = pc.CreateDenseIndexWithIntegratedEmbedding(ctx, IndexName, EmbeddingModel)
	if err != nil {
		switch err {
		case ErrIndexAlreadyExists:
			// index already exists, carry on
		default:
			// some other actual error we haven't anticipated
			fmt.Fprintf(W, "failed to create embedding: %v\n", err)
			os.Exit(1)
		}
	}

	flamehamsterFile, err := os.ReadFile("./ref/flamehamster.md")
	if err != nil {
		fmt.Fprintf(W, "flamehamster open file error: %\n", err)
		os.Exit(1)
	}

	chunkedFile := SplitByH1(string(flamehamsterFile))

	idx, err := pc.client.DescribeIndex(ctx, IndexName)
	handleError(err, "describe index")

	indexConnParams := pinecone.NewIndexConnParams{Host: idx.Host, Namespace: "flamehamster"}

	idxConnection, err := pc.client.Index(indexConnParams)
	handleError(err, "index connection")
	defer idxConnection.Close()

	var records []*pinecone.IntegratedRecord

	for i, chunk := range chunkedFile {
		newRecord := pinecone.IntegratedRecord{
			"id":         fmt.Sprintf("chunk-{%d}", i),
			"chunk_text": chunk,
			"manual":     ManualName,
		}

		records = append(records, &newRecord)
	}

	err = idxConnection.UpsertRecords(ctx, records)
	handleError(err, "upsert records")

	fmt.Printf("Successfully upserted %d record(s)!\n", len(records))
}
