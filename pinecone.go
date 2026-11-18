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

func (c *EmbeddingClient) createDenseIndexWithIntegratedEmbedding(ctx context.Context, indexName, model string) error {
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

func (c *EmbeddingClient) fetchRelevantChunksFromVectorDB(ctx context.Context, idxConnection *pinecone.IndexConnection, userInput string, topk int32) (*pinecone.SearchRecordsResponse, error) {
	res, err := idxConnection.SearchRecords(ctx, &pinecone.SearchRecordsRequest{
		Query: pinecone.SearchRecordsQuery{
			TopK: topk,
			Inputs: &map[string]interface{}{
				"text": userInput,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("search for relevant records: %v", err)
	}

	// // Convert relevant chunks to giant string for prompt injection.
	// relevantText, err := c.convertRelevantHitsToText(res)
	// if err != nil {
	// 	return "", fmt.Errorf("convert relevant hits to text: %w", err)
	// }

	return res, nil
}

// Convert relevant chunks to giant string for prompt injection.
func (c *EmbeddingClient) convertRelevantHitsToText(records *pinecone.SearchRecordsResponse) (string, error) {
	var relevantText string
	for _, hit := range records.Result.Hits {
		chunkText, ok := hit.Fields["chunk_text"]
		if !ok {
			return "", fmt.Errorf("chunk text field doesn't exist: %v", hit)
		}

		text, ok := chunkText.(string)
		if !ok {
			return "", fmt.Errorf("chunk text field is not string: %v", text)
		}

		relevantText += fmt.Sprintf("%s\n\n", text)
	}

	return relevantText, nil
}

func (c *EmbeddingClient) CreateNewNamespace(ctx context.Context, idc *pinecone.IndexConnection, name string) (string, error) {
	ns, err := idc.CreateNamespace(ctx, &pinecone.CreateNamespaceParams{
		Name: name,
	})
	if err != nil {
		return "", fmt.Errorf("cannot create new namespace '%s': %w", name, err)
	}

	return ns.Name, nil
}
