package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/pinecone-io/go-pinecone/v7/pinecone"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	GrossAppPineconeNS string = "all-gross"
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
	ErrIndexAlreadyExists     = errors.New("index already exists")
	ErrNamespaceAlreadyExists = errors.New("namespace already exists")
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

func (c *EmbeddingClient) upsertChunksToIndex(ctx context.Context, idxConnection *pinecone.IndexConnection, manualName string, chunks []string) error {
	var records []*pinecone.IntegratedRecord

	for i, chunk := range chunks {
		newRecord := pinecone.IntegratedRecord{
			"id":         fmt.Sprintf("%s-chunk-%d", manualName, i),
			"chunk_text": chunk,
			"manual":     manualName,
		}

		records = append(records, &newRecord)
	}

	err := idxConnection.UpsertRecords(ctx, records)
	if err != nil {
		return fmt.Errorf("upsert records: %w", err)
	}

	slog.Info("upsert successful", slog.String("manual", manualName), slog.Int("count", len(chunks)))

	return nil
}

func (c *EmbeddingClient) upsertChunkMapToIndex(ctx context.Context, idxConnection *pinecone.IndexConnection, cMap map[string][]string) error {
	for manualName, chunks := range cMap {
		withoutExt := strings.TrimSuffix(manualName, ".md")

		if err := c.upsertChunksToIndex(ctx, idxConnection, withoutExt, chunks); err != nil {
			return fmt.Errorf("upsertChunkMapToIndex loop: %w", err)
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

func (c *EmbeddingClient) createNewNamespace(ctx context.Context, idc *pinecone.IndexConnection, name string) (string, error) {
	ns, err := idc.CreateNamespace(ctx, &pinecone.CreateNamespaceParams{
		Name: name,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return "", fmt.Errorf("standard go error: %w", err)
		}

		switch st.Code() {
		case codes.AlreadyExists:
			return "", ErrNamespaceAlreadyExists
		default:
			return "", fmt.Errorf("gRPC error: [%s]: %w", st.Message(), st.Err())
		}
	}

	return ns.Name, nil
}
