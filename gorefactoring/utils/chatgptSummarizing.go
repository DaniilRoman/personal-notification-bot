package modules

import (
	"context"
	"log"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func ConfigureOpenAI(apiKey string) *openai.Client {
    return openai.NewClient(apiKey)
}

func SummarizeText(text string, client *openai.Client) string {
    modelLimit := 4000
    if len(text) > modelLimit {
        text = text[:4000]
    }
	inputText := "Give me the main ideas of this tech article: " + text
	model := openai.GPT4
	maxTokens := 300

	request := openai.ChatCompletionRequest{
		Model:      model,
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    openai.ChatMessageRoleUser,
                Content:  inputText,
            },
        },
		Temperature: 0,
		MaxTokens:   maxTokens,
	}

	resp, err := client.CreateChatCompletion(
        context.Background(),
        request,
    )
	if err != nil {
		log.Println("Couldn't get summary for the text:", err)
		return ""
	}

	summary := resp.Choices[0].Message.Content
	return strings.TrimSpace(summary)
}