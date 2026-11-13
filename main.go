package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
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

	RelevanceCountTopK int32 = 3
)

func main() {
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

	// client := openai.NewClient(
	// 	option.WithAPIKey(apikey),
	// )

	// chatbot := chatClient{
	// 	llmClient:            &client,
	// 	modelChoice:          openai.ChatModelGPT5_4Mini2026_03_17,
	// 	modelTemperature:     temp,
	// 	modelMaxOutputTokens: MaxOutputTokens,
	// }

	pc, err := NewPineconeClient(pineconeApiKey)
	if err != nil {
		fmt.Fprintf(W, "failed to create pinecone client: %v\n", err)
		os.Exit(1)
	}

	// Create a database instance in Pinecone we can use for semantic search.
	err = pc.CreateDenseIndexWithIntegratedEmbedding(ctx, IndexName, EmbeddingModel)
	if err != nil {
		switch err {
		case ErrIndexAlreadyExists: // index already exists, carry on
		default:
			fmt.Fprintf(W, "failed to create embedding: %v\n", err)
			os.Exit(1)
		}
	}

	// Create a connection to the flamehamster vector database.
	idx, err := pc.client.DescribeIndex(ctx, IndexName)
	handleError(err, "describe index")

	indexConnParams := pinecone.NewIndexConnParams{Host: idx.Host, Namespace: "flamehamster"}

	// Create the actual connection to the flamehamster database on Pinecone.
	idxConnection, err := pc.client.Index(indexConnParams)
	handleError(err, "index connection")
	defer idxConnection.Close()

	// ? Up to here: Start chatbot interaction.

	scanner := bufio.NewScanner(os.Stdin)

	// systemInstruction := "You are an AI customer support technician who is knowledgeable about software products created by the company called GROSS. One such product is a web browser called Flamehamster."
	initialAssistantMessage := "How can I help you today?"

	fmt.Printf(TextColourAssistant+"Assistant: %s\n", initialAssistantMessage)

	fmt.Print(TextColourUser + "User: ")
	scanner.Scan()
	userInput := scanner.Text()

	// conversationHistory += ChatbotHelpfulAssistantInitial + userInput
	// conversationHistory := []responses.ResponseInputItemUnionParam{
	// 	newSystemInstruction(systemInstruction),
	// 	newAssistantMessage(initialAssistantMessage),
	// }

	// Chat loop:
	for userInput != "exit" {

		// ! Working here:
		// Process relevant results from vector db:
		res, err := pc.FetchRelevantChunksFromVectorDB(ctx, idxConnection, userInput, RelevanceCountTopK)
		handleError(err, "retrieving relevant chunks from pinecone")

		var relevantText string
		for _, hit := range res.Result.Hits {
			chunkText, ok := hit.Fields["chunk_text"]
			if !ok {
				fmt.Fprintf(W, "chunk text field doesn't exist: %v", hit)
				os.Exit(1)
			}

			text, ok := chunkText.(string)
			if !ok {
				fmt.Fprintf(W, "chunk text field is not string: %v", text)
				os.Exit(1)
			}

			relevantText += fmt.Sprintf("%s\n\n", text)
		}

		// ? Pausing on chatbot interaction for a sec... just to test it's working.
		// // Chat interaction:
		// response, err := chatbot.sendMessageWithHistory(ctx, conversationHistory)
		// if err != nil {
		// 	fmt.Fprintf(W, "generation failed: %v", err)
		// 	log.Println(conversationHistory)
		// 	os.Exit(1)
		// }

		// llmResponseText := fmt.Sprintf(TextColourAssistant+"Assistant: %s\n", response.OutputText())
		// fmt.Fprint(W, llmResponseText)

		fmt.Fprintf(W, TextColourAssistant+"Relevant text chunks:\n%s\n", relevantText)

		fmt.Fprint(W, TextColourUser+"User: ")
		scanner.Scan()
		userInput = scanner.Text()

		// // Append LLM output response to history.
		// conversationHistory = append(
		// 	conversationHistory,
		// 	outputAsInput(response.Output)...,
		// )

		// // Append user input prompt to history.
		// conversationHistory = append(
		// 	conversationHistory,
		// 	newUserMessage(userInput),
		// )
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(W, "scanner error: %v", err)
		os.Exit(1)
	}
}
