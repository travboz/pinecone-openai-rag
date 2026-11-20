package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/pinecone-io/go-pinecone/v7/pinecone"
)

func RunEvaluationPipeline(ctx context.Context, chatClient *chatClient, pc *EmbeddingClient, idxConn *pinecone.IndexConnection) {
	inputFile, err := os.Open("./evaluations/queries.csv")
	handleError(err, "opening input file")
	defer inputFile.Close()

	outputFile, err := os.Create("./evaluations/traces.csv")
	handleError(err, "creating output file")
	defer outputFile.Close()

	csvReader := csv.NewReader(inputFile)
	csvWriter := csv.NewWriter(outputFile)

	err = csvWriter.Write([]string{"Query Topic", "User Query", "History", "AI Response"})
	handleError(err, "writing header row to file")

	for {
		// Resetting history:
		var history conversationHistory

		history.AddSystemInstruction(SystemInstructions)
		history.AddAssistantMessage(InitialAssistantMessage)

		row, err := csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				handleError(err, "reading from input csv")
			}
		}

		queryTopic := row[0]
		userInput := row[1]

		// Find and process relevant results from vector db:
		relevantKnowledge, err := pc.FetchRelevantKnowledge(ctx, idxConn, userInput)
		handleError(err, "retrieving relevant chunks from pinecone")

		// Augment user input prompt with retrieved, relevant text chunks:
		augmentedPrompt := augmentUserInputWithRelevantKnowledge(userInput, relevantKnowledge, RetrievalAugmentPrompt)

		// * Append augmented user prompt to history.
		history.AddUserMessage(augmentedPrompt)

		// Send augmented prompt off to LLM for a response.
		response, err := chatClient.sendMessageWithHistory(ctx, history)
		handleError(err, "getting response from LLM")

		historyBytes, err := json.MarshalIndent(history, "", "  ")
		handleError(err, "marshaling history to json")

		err = csvWriter.Write([]string{queryTopic, userInput, string(historyBytes), response.OutputText()})
		handleError(err, "writing to output csv")

		csvWriter.Flush()

		slog.Info(
			"completed trace",
			slog.String("query_topic", queryTopic),
			slog.String("user_input", userInput),
		)
	}
}
