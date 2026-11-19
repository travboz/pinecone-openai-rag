package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/openai/openai-go/v3"
	"github.com/pinecone-io/go-pinecone/v7/pinecone"
	"github.com/spf13/pflag"
)

var (
	TextColourUser      = "\033[34m"
	TextColourAssistant = "\033[32m"
)

var (
	MaxOutputTokens int64     = 1024
	W               io.Writer = os.Stdout

	EmbeddingModel     string  = "llama-text-embed-v2"
	IndexName          string  = "gross-app"
	temp               float64 = 0
	RelevanceCountTopK int32   = 3
	CreateRecords      bool
)

func main() {
	pflag.BoolVar(&CreateRecords, "build", false, "builds the index, namespace, and upserts documents for project")
	pflag.Parse()

	// Retrieve relevant API keys for services.
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

	// Create ChatGPT and Pinecone clients.
	ctx := context.TODO()

	// ? Create connection to Pinecone and set up knowledge base.
	pc, err := NewPineconeClient(pineconeApiKey)
	if err != nil {
		fmt.Fprintf(W, "failed to create pinecone client: %v\n", err)
		os.Exit(1)
	}

	// Create a database instance in Pinecone we can use for semantic search.
	err = pc.createDenseIndexWithIntegratedEmbedding(ctx, IndexName, EmbeddingModel)
	if err != nil {
		switch err {
		case ErrIndexAlreadyExists: // index already exists, carry on
			slog.Info("no operation", slog.String("reason", err.Error()))
		default:
			fmt.Fprintf(W, "failed to create embedding: %v\n", err)
			os.Exit(1)
		}
	}

	// Create a connection to the flamehamster vector database.
	idx, err := pc.client.DescribeIndex(ctx, IndexName)
	handleError(err, "describe index")

	indexConnParams := pinecone.NewIndexConnParams{Host: idx.Host, Namespace: "all-gross"}

	// Create the actual connection to the flamehamster database on Pinecone.
	idxConnection, err := pc.client.Index(indexConnParams)
	handleError(err, "index connection")
	defer idxConnection.Close()

	if CreateRecords {
		// Upsert all document chunks into Pinecone using the new namespace.
		// Create chunks for each file.
		chunksMap, err := ChunksFromDirectory("./data")
		handleError(err, "unable to chunk files in directory")

		// Iterate over chunks and upsert into Pinecone.
		err = pc.upsertChunkMapToIndex(ctx, idxConnection, chunksMap)
		handleError(err, "unable to upsert chunks to index")
	}

	// Chat client creation.
	chatClient := NewChatClient(
		apikey,
		openai.ChatModelGPT5_4Mini2026_03_17,
		temp,
		MaxOutputTokens,
	)

	// * Start chatbot interaction.
	scanner := bufio.NewScanner(os.Stdin)

	// Creating initial history:
	var history conversationHistory

	history.AddSystemInstruction(SystemInstructions)
	history.AddAssistantMessage(InitialAssistantMessage)

	// Printing 'How can I help you':
	fmt.Printf(TextColourAssistant+"Assistant: %s\n", InitialAssistantMessage)

	// * Chat loop:
	for {
		fmt.Fprint(W, TextColourUser+"User: ")

		if !scanner.Scan() {
			fmt.Fprintln(W) // EOF (Ctrl+D) or read error
			break
		}

		userInput := scanner.Text()
		if userInput == "exit" {
			break
		}

		if userInput == "" {
			continue // don't spend a request on an empty turn
		}

		// Find and process relevant results from vector db:
		relevantKnowledge, err := pc.FetchRelevantKnowledge(ctx, idxConnection, userInput)
		handleError(err, "retrieving relevant chunks from pinecone")

		slog.Info("relevant knowledge retrieved", slog.String("knowledge", relevantKnowledge))

		// Augment user input prompt with retrieved, relevant text chunks:
		augmentedPrompt := augmentUserInputWithRelevantKnowledge(userInput, relevantKnowledge, RetrievalAugmentPrompt)

		// * Append augmented user prompt to history.
		history.AddUserMessage(augmentedPrompt)

		// Send augmented prompt off to LLM for a response.
		response, err := chatClient.sendMessageWithHistory(ctx, history)
		if err != nil {
			slog.Error("generation failed", slog.Any("error", err), slog.Any("history", history))
			os.Exit(1)
		}

		history.AddLLMResponse(response)
		fmt.Fprintf(W, TextColourAssistant+"Assistant: %s\n", response.OutputText())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(W, "scanner error: %v", err)
		os.Exit(1)
	}
}
