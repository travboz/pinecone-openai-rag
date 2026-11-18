package main

import (
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
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

	// ? Todo: Upsert all files into Pinecone
	// Create chunks for each file.
	chunksMap, err := ChunksFromDirectory("./data")
	handleError(err, "unable to chunk files in directory")

	// Iterate over chunks and upsert into Pinecone.
	err = pc.upsertChunkMapToIndex(chunksMap)
	handleError(err, "unable to upsert chunks to index")

	/* // ! Delete this comment opener to continue
	// Create ChatGPT and Pinecone clients.
	ctx := context.TODO()

	chatClient := NewChatClient(
		apikey,
		openai.ChatModelGPT5_4Mini2026_03_17,
		temp,
		MaxOutputTokens,
	)

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

	systemInstruction := "You are an AI customer support technician who is knowledgeable about software products created by the company called GROSS. One such product is a web browser called Flamehamster."
	initialAssistantMessage := "How can I help you today?"
	// Order is {relevant_docs} then {original_user_input}.
	retrievalAugmentPrompt := "Here are excerpts from the official Flamehamster web browser documentation: %s. Use whatever info from the above documentation excerpts (and no other info) to answer the following query: %s"

	// Printing How can I help you:
	fmt.Printf(TextColourAssistant+"Assistant: %s\n", initialAssistantMessage)

	// Getting first user input:
	fmt.Print(TextColourUser + "User: ")
	scanner.Scan()
	userInput := scanner.Text()

	// Creating initial history:
	conversationHistory := []responses.ResponseInputItemUnionParam{
		newSystemInstruction(systemInstruction),
		newAssistantMessage(initialAssistantMessage),
	}

	// * Chat loop:
	for userInput != "exit" {

		// Process relevant results from vector db:
		res, err := pc.fetchRelevantChunksFromVectorDB(ctx, idxConnection, userInput, RelevanceCountTopK)
		handleError(err, "retrieving relevant chunks from pinecone")

		relevantText, err := pc.convertRelevantHitsToText(res)
		handleError(err, "converting relevant results to text")

		augmentedPrompt := augmentUserInputWithRelevantKnowledge(userInput, relevantText, retrievalAugmentPrompt)

		// Append augmented user input prompt to history.
		conversationHistory = append(conversationHistory, newUserMessage(augmentedPrompt))

		// * Chat interaction:
		// Augment user input prompt with retrieved, relevant text chunks:
		response, err := chatClient.sendMessageWithHistory(ctx, conversationHistory)
		if err != nil {
			fmt.Fprintf(W, "generation failed: %v", err)
			log.Println(conversationHistory)
			os.Exit(1)
		}

		llmResponseText := fmt.Sprintf(TextColourAssistant+"Assistant: %s\n", response.OutputText())
		fmt.Fprint(W, llmResponseText)

		// Append LLM output response to history.
		conversationHistory = append(
			conversationHistory,
			outputAsInput(response.Output)...,
		)

		fmt.Fprint(W, TextColourUser+"User: ")
		scanner.Scan()
		userInput = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(W, "scanner error: %v", err)
		os.Exit(1)
	}
	*/
}
