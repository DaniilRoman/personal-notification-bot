package main

import (
	blogs "main/modules/blogs"
	blogsStats "main/modules/blogsStats"
	currency "main/modules/currency"
	hertha "main/modules/herthaTickets"
	mobilenumber "main/modules/mobileNumber"
	union "main/modules/unionBerlinTickets"
	weather "main/modules/weather"
	word "main/modules/wordOfTheDay"
	justAiNews "main/modules/justAiNews"
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


func main() {
	dynamodb := utils.NewDynamoDbService(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, REGION_NAME, nil)
	chatGptService := utils.NewChatGptService(OPENAI_ACCESS_KEY)

	var wg sync.WaitGroup
	wg.Add(8)
	weatherChan := make(chan *weather.WeatherData, 1)
	currencyChan := make(chan *currency.CurrencyData, 1)
	wordOfTheDayChan := make(chan *word.WordOfTheDayData, 1)
	herthaTicketsChan := make(chan *hertha.HerthaTicketsData, 1)
	unionBerlinTicketsChan := make(chan *union.UnionBerlinTicketsData, 1)
	mobileNumberChan := make(chan *mobilenumber.MobileNumberData, 1)
	blogsChan := make(chan *blogs.BlogUpdateData, 1)
	justAiNewsChan := make(chan *justAiNews.JustAiNewsData, 1)

	go func() {
	    weatherChan <- weather.GetWeather(OPEN_WHEATHER_API_KEY)
		wg.Done()	
	}()

	go func() {
	    currencyChan <- currency.Currency(EXCHANGERATE_API_KEY)
		wg.Done()	
	}()

	go func() {
	    wordOfTheDayChan <- word.WordOfTheDay()
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
	    justAiNewsChan <- justAiNews.JustAiNews(dynamodb)
		wg.Done()	
	}()

	wg.Wait()

	wordOfTheDayData := <- wordOfTheDayChan
	herthaTicketsData := <- herthaTicketsChan
	unionBerlinTicketsData := <- unionBerlinTicketsChan
	weatherData := <- weatherChan
	currencyData := <- currencyChan
	blogsUpdatesData := <- blogsChan
	mobileNimberData := <- mobileNumberChan
	justAiNewsData := <- justAiNewsChan

	blogsStatsData := blogsStats.BlogsStats(blogsUpdatesData.PopularWords(), dynamodb, chatGptService)


	utils.SendToTelegram(TELEGRAM_TOKEN, TELEGRAM_CHAT_ID,
		blogsUpdatesData,
	)

	utils.SendToTelegram(TELEGRAM_TOKEN, TELEGRAM_CHAT_ID,
		weatherData,
		currencyData,
		herthaTicketsData,
		unionBerlinTicketsData,
		mobileNimberData,
		wordOfTheDayData,
		justAiNewsData,
		blogsStatsData,
	)


	dataForRendering := map[string]interface{} {
		"Weather" : weatherData,
		"Currency" : currencyData,
		"Blogs" : blogsUpdatesData,
	}
	utils.RenderIndexHTML(dataForRendering)
}
