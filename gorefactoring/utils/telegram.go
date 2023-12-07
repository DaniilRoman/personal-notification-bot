package modules

import (
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendToTelegram(token string, chatId int64, data ...toString) {
	message := telegramData(data...)

	bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        log.Fatal("Couldn't initialise Telegram bot Api", err)
		return
    }
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ParseMode = "Markdown"
	if _, err := bot.Send(msg); err != nil {
        log.Fatal("Couldn't send a message to Telegram", err)
	}
}

func telegramData(data ...toString) string {
	res := []string {}
	for _, el := range data {
		if el != nil {
			str := el.String()
			if str != "" {
				res = append(res, str)
			}
	
		}
	}
    res = append(res, "[Html page](https://daniilroman.github.io/personal-notification-bot/)")
	return strings.Join(res, "\n\n")
}

type toString interface {
	String() string
}
