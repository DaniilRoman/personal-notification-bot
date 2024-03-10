package blogsStats

import (
	"log"
	"main/utils"
	"time"
)

var today time.Time = time.Now()
var dayFormat string = "2006-01-02"
var monthFormat string = "2006-01"


func BlogsStats(popularWords string, dynamodb *utils.DynamoDbService, chatGpt *utils.ChatGptService) *BlogsStatsData {
    saveTodaysStats(popularWords, dynamodb)

    sundayStats := ""
    monthStats := ""

    if sunday() {
		log.Printf("Runnung a statistic for a week...")
        collectedPopularWords := popularWordsFromPrevWeek(dynamodb)
        collectedPopularWords += ","+popularWords        
        weekPopularWords := chatGpt.AggregatedPopularWords(collectedPopularWords)
        saveWeekStats(weekPopularWords, dynamodb)
        sundayStats = weekPopularWords
		log.Printf("Finished a statistic for a week.")
    }

    if lastDayOfMonth() {
		log.Printf("Runnung a statistic for a month...")
        collectedPopularWords := dynamodb.GetBlogsStat(today.Format(monthFormat))
        monthStats = chatGpt.AggregatedPopularWords(collectedPopularWords)
		log.Printf("Finished a statistic for a month.")
    }

    return &BlogsStatsData{sundayStats, monthStats}
}

func sunday() bool {
    weekday := today.Weekday()
    return weekday == time.Sunday 
}

func lastDayOfMonth() bool {
    todayMonth := today.Format(monthFormat)
    tomorrowMonth := today.AddDate(0, 0, 1).Format(monthFormat)
    return todayMonth != tomorrowMonth 
}

func saveTodaysStats(popularWords string, dynamodb *utils.DynamoDbService) {
    if popularWords == "" {
        return
    }
    itemKey := today.Format(dayFormat)
    dynamodb.SavePopularWords(itemKey, popularWords)
}

func saveWeekStats(popularWords string, dynamodb *utils.DynamoDbService) {
    itemKey := today.Format(monthFormat)
    dynamodb.AppendPopularWords(itemKey, popularWords)
}

func popularWordsFromPrevWeek(dynamodb *utils.DynamoDbService) string {
    day1 := today.AddDate(0, 0, -1).Format(dayFormat)
    day2 := today.AddDate(0, 0, -2).Format(dayFormat)
    day3 := today.AddDate(0, 0, -3).Format(dayFormat)
    day4 := today.AddDate(0, 0, -4).Format(dayFormat)
    day5 := today.AddDate(0, 0, -5).Format(dayFormat)
    day6 := today.AddDate(0, 0, -6).Format(dayFormat)
    return dynamodb.GetStatsFromPrevDays([]string{day1, day2, day3, day4, day5, day6})
}
