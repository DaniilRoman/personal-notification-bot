package utils

import (
	"context"
	"log"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type ChatGptService struct {
	client *openai.Client
}

func NewChatGptService(apiKey string) *ChatGptService {
	client := configureOpenAI(apiKey)
	return &ChatGptService{client}
}

func configureOpenAI(apiKey string) *openai.Client {
    return openai.NewClient(apiKey)
}

func (service *ChatGptService) SummarizeText(text string) string {
	prompt := "Why would I want to read this tech article?\n"
	return service.chatCompletionRequest(text, prompt)
}

func (service *ChatGptService) ArticlePopularWords(text string) string {
	prompt := "Here's a tech article below. I'm a developer and my goal is to understand the main idea of this article from technologies point of view. Can you show me the top 10 most popular words that are related to technologies in this tech article? But show it without any formatting, just separate by comma.\n"
	return service.chatCompletionRequest(text, prompt)
}

func (service *ChatGptService) AggregatedPopularWords(text string) string {
	prompt := "Here're the keywords from the different technical articles. They might diverse but have the same meaning at the same time. Can you show me the top 10 most popular words? But show it without any formatting, just separate by comma.\n"
	return service.chatCompletionRequest(text, prompt)
}

func (service *ChatGptService) CommonSummaryFromSeveralNews(text string) string {
	prompt := "Here're the several different technical article summaries separated by a new line. Could you please summarise what happend in these articles.\n"
	return service.chatCompletionRequest(text, prompt)
}

func (service *ChatGptService) chatCompletionRequest(text string, prompt string) string {
    modelLimit := 16000
    if len(text) > modelLimit {
        text = text[:modelLimit]
    }
	inputText := prompt + text
	model := openai.GPT3Dot5Turbo
	maxTokens := 200

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

	resp, err := service.client.CreateChatCompletion(
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
