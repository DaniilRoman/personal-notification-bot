package main

import (
	blogs "main/modules/blogs"
	currency "main/modules/currency"
	expenses "main/modules/expenses"
	gmailreader "main/modules/gmailReader"
	hertha "main/modules/herthaTickets"
	justAiNews "main/modules/justAiNews"
	mobilenumber "main/modules/mobileNumber"
	simpleanalytics "main/modules/simpleanalytics"
	union "main/modules/unionBerlinTickets"
	weather "main/modules/weather"
	word "main/modules/wordOfTheDay"
	"main/utils"
	"os"
	"strconv"
	"sync"
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

var APP_SCRIPT_ID = os.Getenv("APP_SCRIPT_ID")

type separatedSender interface {
	GetUpdateStrings() []string
}

func sendSeparately(data separatedSender) {
	updateStrings := data.GetUpdateStrings()
	if len(updateStrings) == 0 {
		return
	}
	for _, updateStr := range updateStrings {
		utils.SendToTelegram(TELEGRAM_TOKEN, TELEGRAM_CHAT_ID, updateStr)
	}
}

func main() {
	dynamodb := utils.NewDynamoDbService(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, REGION_NAME, nil)
	chatGptService := utils.NewChatGptService(OPENAI_ACCESS_KEY)

	var wg sync.WaitGroup
	wg.Add(11)
	weatherChan := make(chan *weather.WeatherData, 1)
	currencyChan := make(chan *currency.CurrencyData, 1)
	wordOfTheDayChan := make(chan *word.WordOfTheDayData, 1)
	herthaTicketsChan := make(chan *hertha.HerthaTicketsData, 1)
	unionBerlinTicketsChan := make(chan *union.UnionBerlinTicketsData, 1)
	mobileNumberChan := make(chan *mobilenumber.MobileNumberData, 1)
	blogsChan := make(chan *blogs.BlogUpdateData, 1)
	gmailReaderChan := make(chan *gmailreader.GmailReaderData, 1)
	justAiNewsChan := make(chan *justAiNews.JustAiNewsData, 1)
	expensesTotalChan := make(chan *expenses.MonthlyTotalData, 1)

	go func() {
		weatherChan <- weather.GetWeather(OPEN_WHEATHER_API_KEY)
		wg.Done()
	}()

	go func() {
		currencyChan <- currency.Currency(EXCHANGERATE_API_KEY)
		wg.Done()
	}()

	go func() {
		wordOfTheDayChan <- word.WordOfTheDay(dynamodb)
		wg.Done()
	}()

	go func() {
		herthaTicketsChan <- hertha.HerthaTickets(dynamodb)
		wg.Done()
	}()

	go func() {
		unionBerlinTicketsChan <- union.UnionBerlinTickets(dynamodb)
		wg.Done()
	}()

	go func() {
		mobileNumberChan <- mobilenumber.MobileNumberNotification()
		wg.Done()
	}()

	go func() {
		blogsChan <- blogs.BlogUpdates(chatGptService)
		wg.Done()
	}()

	go func() {
		gmailReaderChan <- gmailreader.GmailReader()
		wg.Done()
	}()

	go func() {
		justAiNewsChan <- justAiNews.JustAiNews(dynamodb)
		wg.Done()
	}()

	go func() {
		simpleanalytics.MakeSimpleAnalyticsScreenshot()
		wg.Done()
	}()

	go func() {
		expensesTotalChan <- expenses.GetMonthlyTotal()
		wg.Done()
	}()

	wg.Wait()

	wordOfTheDayData := <-wordOfTheDayChan
	herthaTicketsData := <-herthaTicketsChan
	unionBerlinTicketsData := <-unionBerlinTicketsChan
	weatherData := <-weatherChan
	currencyData := <-currencyChan
	blogsUpdatesData := <-blogsChan
	gmailReaderData := <-gmailReaderChan
	mobileNimberData := <-mobileNumberChan
	justAiNewsData := <-justAiNewsChan
	expensesTotal := <-expensesTotalChan

	utils.SendImagesToTelegram(TELEGRAM_TOKEN, TELEGRAM_CHAT_ID)

	for _, blog := range blogsUpdatesData.Blogs {
		utils.SendToTelegramWithButton(TELEGRAM_TOKEN, TELEGRAM_CHAT_ID, blog.String(), "Summarize", "summarize")
	}
	sendSeparately(gmailReaderData)

	utils.SendToTelegramWithInterface(TELEGRAM_TOKEN, TELEGRAM_CHAT_ID,
		weatherData,
		currencyData,
		herthaTicketsData,
		unionBerlinTicketsData,
		mobileNimberData,
		wordOfTheDayData,
		justAiNewsData,
		expensesTotal,
	)

	dataForRendering := map[string]interface{}{
		"Weather":     weatherData,
		"Currency":    currencyData,
		"Blogs":       blogsUpdatesData,
		"Expenses":    expensesTotal,
		"AppScriptId": APP_SCRIPT_ID,
	}
	utils.RenderWwwResources(dataForRendering)
}
