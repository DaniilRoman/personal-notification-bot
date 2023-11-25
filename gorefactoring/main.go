package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"main/modules"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


var OPEN_WHEATHER_API_KEY = os.Getenv("OPEN_WHEATHER_API_KEY")
var EXCHANGERATE_API_KEY = os.Getenv("EXCHANGERATE_API_KEY")

var TELEGRAM_TOKEN = os.Getenv("TELEGRAM_TOKEN")
var TELEGRAM_CHAT_ID, _ = strconv.ParseInt(os.Getenv("TELEGRAM_TO"), 10, 64)

var AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
var AWS_SECRET_ACCESS_KEY = os.Getenv("AWS_SECRET_ACCESS_KEY")
var REGION_NAME = os.Getenv("REGION_NAME")

var OPENAI_ACCESS_KEY = os.Getenv("OPENAI_ACCESS_KEY")
var OPENAI_ORGANIZATION = os.Getenv("OPENAI_ORGANIZATION")

func main() {
    currencyData 			:= modules.Currency(EXCHANGERATE_API_KEY)
	weatherData 			:= modules.GetWeather(OPEN_WHEATHER_API_KEY)
	blogsUpdatesData 		:= modules.BlogUpdates()
	herthaTicketsData 		:= modules.HerthaTickets()
	unionBerlinTicketsData  := modules.UnionBerlinTickets()
	mobileNimberData 		:= modules.MobileNumberNotification()
	wordOfTheDay   			:= modules.WordOfTheDay()

	telegramData := telegramData(
		currencyData,
		weatherData,
		blogsUpdatesData,
		herthaTicketsData,
		unionBerlinTicketsData,
		mobileNimberData,
		wordOfTheDay,
	)

	sendToTelegram(TELEGRAM_TOKEN, TELEGRAM_CHAT_ID, telegramData)

	// dataForRendering := map[string]interface{} {
	// 	"Weather" : weatherData,
	// 	"Currency" : currencyData,
	// 	"HerthaTickets" : herthaTicketsData,
	// 	"Blogs" : blogsUpdatesData,
	// }
	// utils.RenderIndexHTML(dataForRendering) 
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

func telegramData(data ...toString) string {
	res := []string {}
	for _, el := range data {
		if el != nil {
			res = append(res, el.String())
		}
	}
	return strings.Join(res, "\n")
}

type toString interface {
	String() string
}