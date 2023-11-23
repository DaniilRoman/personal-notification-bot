package main

import (
	"fmt"
	"log"
	modules "main/modules"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var EXCHANGERATE_API_KEY = os.Getenv("EXCHANGERATE_API_KEY")
var TELEGRAM_TOKEN = os.Getenv("TELEGRAM_APITOKEN")
var TELEGRAM_CHAT_ID, chatIdError = strconv.ParseInt(os.Getenv("TELEGRAM_APITOKEN"), 10, 64)

func main() {
	if chatIdError != nil {
		panic(chatIdError)
	}


    currencyData 			:= modules.Currency(EXCHANGERATE_API_KEY)
	blogsUpdatesData 		:= modules.BlogUpdates()
	herthaTicketsData 		:= modules.HerthaTickets()
	unionBerlinTicketsData  := modules.UnionBerlinTickets()
	mobileNimberData 		:= modules.MobileNumberNotification()

	telegramData := telegramData(
		[]toString{
			currencyData,
			blogsUpdatesData,
			herthaTicketsData,
			unionBerlinTicketsData,
			mobileNimberData,
		},
	)
	sendToTelegram(TELEGRAM_TOKEN, TELEGRAM_CHAT_ID, telegramData)

	fmt.Print("=== END ===")
}

func sendToTelegram(token string, chatId int64, message string) {
	bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        log.Fatal("Couldn't initialise Telegram bot Api", err)
		return
    }
	msg := tgbotapi.NewMessage(chatId, message)
	if _, err := bot.Send(msg); err != nil {
        log.Fatal("Couldn't send a message to Telegram", err)
	}

}

func telegramData(data []toString) string {
	res := []string {}
	for _, el := range data {
		res = append(res, el.String())
	}
	return strings.Join(res, "\n")
}

type toString interface {
	String() string
}