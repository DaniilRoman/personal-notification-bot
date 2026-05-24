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
	// TODO either refactor and fix or remove completely
	// prompt := "Here's a tech article below. I'm a developer and my goal is to understand the main idea of this article from technologies point of view. Can you show me the top 10 most popular words that are related to technologies in this tech article? But show it without any formatting, just separate by comma.\n"
	// return service.chatCompletionRequest(text, prompt)
	return ""
}

func (service *ChatGptService) AggregatedPopularWords(text string) string {
	prompt := "Here're the keywords from the different technical articles. They might diverse but have the same meaning at the same time. Can you show me the top 10 most popular words? But show it without any formatting, just separate by comma.\n"
	return service.chatCompletionRequest(text, prompt)
}

func (service *ChatGptService) CommonSummaryFromSeveralNews(text string) string {
	prompt := "Here're the several different technical article summaries separated by a new line. Could you please summarise what happend in these articles.\n"
	return service.chatCompletionRequest(text, prompt)
}

func (service *ChatGptService) BlogsPodcastSummary(text string) string {
	model := "gpt-5.4-mini"
	prompt := `Ты готовишь утренний аудио-брифинг, который будет зачитываться вслух умным домашним помощником. 
		Я буду слушать его как свою ежедневную утреннюю сводку о технологиях — как короткий сегмент техно-подкаста. 
		Ниже приведён список сегодняшних постов из технологических блогов с их заголовками, описаниями и ссылками.

		Напиши естественное, разговорное устное резюме в стиле дружелюбного ведущего утренних новостей. 
		Текст должен занимать примерно 2–3 минуты при чтении вслух. 
		НЕ используй URL-адреса, markdown, списки, заголовки или какое-либо форматирование — только обычные разговорные предложения. 
		Где это уместно, объединяй связанные темы. 
		Сосредоточься на самых интересных и значимых новостях.

		ВАЖНО: отвечай только на русском языке.

		Сегодняшние посты из блогов:`

	modelLimit := 16000
	if len(text) > modelLimit {
		text = text[:modelLimit]
	}
	inputText := prompt + text

	request := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: inputText,
			},
		},

		Temperature:         0.7,
		MaxCompletionTokens: 3000,
	}

	resp, err := service.client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		log.Println("Couldn't generate podcast summary:", err)
		return ""
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content)
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
