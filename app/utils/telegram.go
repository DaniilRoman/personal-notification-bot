package utils

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendImagesToTelegram(token string, chatId int64) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	files, err := os.ReadDir("./data/images")
	if err != nil {
		return err
	}

	var mediaGroup []interface{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := "./data/images/" + file.Name()
		photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath(filePath))
		mediaGroup = append(mediaGroup, photo)
	}

	if len(mediaGroup) > 0 {
		mediaMsg := tgbotapi.NewMediaGroup(chatId, mediaGroup)
		_, err = bot.Send(mediaMsg)
		if err != nil {
			return err
		}
	}

	return nil
}

func SendToTelegram(token string, chatId int64, data ...string) {
	message := strings.Join(data, "\n\n")
	if message == "" {
		return
	}

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

func SendToTelegramWithInterface(token string, chatId int64, data ...toString) {
	message := telegramData(data...)
	if message == "" {
		return
	}
	SendToTelegram(token, chatId, message)
}

func telegramData(data ...toString) string {
	res := []string{}
	for _, el := range data {
		if el != nil {
			str := el.String()
			if str != "" {
				res = append(res, str)
			}

		}
	}
	return strings.Join(res, "\n\n")
}

type toString interface {
	String() string
}
