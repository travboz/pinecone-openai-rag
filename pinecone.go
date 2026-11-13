package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/pinecone-io/go-pinecone/v7/pinecone"
)

type EmbeddingClient struct {
	client *pinecone.Client
}

func NewPineconeClient(apikey string) (*EmbeddingClient, error) {
	pc, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: apikey,
	})
	if err != nil {
		return nil, fmt.Errorf("create pinecone client: %w", err)
	}

	return &EmbeddingClient{pc}, nil
}

var (
	ErrIndexAlreadyExists = errors.New("index already exists")
)

func (c *EmbeddingClient) CreateDenseIndexWithIntegratedEmbedding(ctx context.Context, indexName, model string) error {
	_, err := c.client.CreateIndexForModel(ctx, &pinecone.CreateIndexForModelRequest{
		Name:   indexName,
		Cloud:  pinecone.CloudAWS,
		Region: "us-east-1",
		Embed: pinecone.CreateIndexForModelEmbed{
			Model:    model,
			FieldMap: map[string]any{"text": "chunk_text"},
		},
	})
	if err != nil {
		var pErr *pinecone.PineconeError
		if errors.As(err, &pErr) && pErr.Code == http.StatusConflict {
			// index already exists
			return ErrIndexAlreadyExists
		} else {
			return fmt.Errorf("create dense index with integrated embedding: %w", err)
		}
	}

	return nil
}
