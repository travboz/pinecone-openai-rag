package main

var (
	ChatbotHelpfulAssistantInitial = "Assistant: How can I help you today?\n\nUser: "

	ChatbotInputPrompt = `I'm a helpful chatbot! Ask me anything:`

	ChatbotPiratePrompt         = `Ahoy! Got questions? Spit 'em out, ye scallywag!`
	ChatbotSystemPirateTemplate = `Respond to the following like a pirate: %s`

	ChatbotSpanishPrompt                         = `Enter a phrase, and I'll translate it into Spanish!`
	ChatbotSystemSpanishTranslateTemplate        = `Translate the following phrase into Spanish: %s`
	ChatbotSystemSpanishTranslateFocusedTemplate = `Translate the following phrase into Spanish: %s. 
	Answer with just the translation and no other introductory text.`
)
