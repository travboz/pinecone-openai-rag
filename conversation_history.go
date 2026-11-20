package main

import "github.com/openai/openai-go/v3/responses"

type conversationHistory []responses.ResponseInputItemUnionParam

func (h *conversationHistory) AddSystemInstruction(inst string) {
	*h = append(*h, newSystemInstruction(inst))
}

func (h *conversationHistory) AddAssistantMessage(msg string) {
	*h = append(*h, newAssistantMessage(msg))
}

func (h *conversationHistory) AddUserMessage(msg string) {
	*h = append(*h, newUserMessage(msg))
}

func (h *conversationHistory) AddLLMResponse(res *responses.Response) {
	// Append LLM output response to history.
	*h = append(
		*h,
		outputAsInput(res.Output)...,
	)
}
