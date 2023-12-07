package main

import (
	"os"
	"strconv"
	"sync"

	"main/modules"
	utils "main/utils"
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
	var wg sync.WaitGroup
	wg.Add(7)
	weatherChan := make(chan *modules.WeatherData, 1)
	currencyChan := make(chan *modules.CurrencyData, 1)
	wordOfTheDayChan := make(chan *modules.WordOfTheDayData, 1)
	herthaTicketsChan := make(chan *modules.HerthaTicketsData, 1)
	unionBerlinTicketsChan := make(chan *modules.UnionBerlinTicketsData, 1)
	mobileNumberChan := make(chan *modules.MobileNumberData, 1)
	blogsChan := make(chan *modules.BlogUpdateData, 1)

	go func() {
	    weatherChan <- modules.GetWeather(OPEN_WHEATHER_API_KEY)
		wg.Done()	
	}()

	go func() {
	    currencyChan <- modules.Currency(EXCHANGERATE_API_KEY)
		wg.Done()	
	}()

	go func() {
	    wordOfTheDayChan <- modules.WordOfTheDay()
		wg.Done()	
	}()

	go func() {
	    herthaTicketsChan <- modules.HerthaTickets()
		wg.Done()	
	}()

	go func() {
	    unionBerlinTicketsChan <- modules.UnionBerlinTickets()
		wg.Done()	
	}()

	go func() {
	    mobileNumberChan <- modules.MobileNumberNotification()
		wg.Done()	
	}()

	go func() {
		client := utils.ConfigureOpenAI(OPENAI_ACCESS_KEY)
	    blogsChan <- modules.BlogUpdates(client)
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


	utils.SendToTelegram(TELEGRAM_TOKEN, TELEGRAM_CHAT_ID,
		weatherData,
		currencyData,
		blogsUpdatesData,
		herthaTicketsData,
		unionBerlinTicketsData,
		mobileNimberData,
		wordOfTheDayData,
	)


	dataForRendering := map[string]interface{} {
		"Weather" : weatherData,
		"Currency" : currencyData,
		"HerthaTickets" : herthaTicketsData,
		"Blogs" : blogsUpdatesData,
	}
	utils.RenderIndexHTML(dataForRendering)
}
