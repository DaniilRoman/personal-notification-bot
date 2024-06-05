package blogsStats

import (
	"log"
	"main/utils"
	"strings"
	"time"
)

var today time.Time = time.Now()
var dayFormat string = "2006-01-02"
var monthFormat string = "2006-01"


func BlogsStats(popularWords string, dynamodb *utils.DynamoDbService, chatGpt *utils.ChatGptService) *BlogsStatsData {
    saveTodaysStats(popularWords, dynamodb)

    sundayStats := ""
    newWeekWords := ""
    monthStats := ""

    if sunday() {
		log.Printf("Runnung a statistic for a week...")
        collectedPopularWords := popularWordsFromPrevWeek(dynamodb)
        collectedPopularWords += ","+popularWords        
        weekPopularWords := chatGpt.AggregatedPopularWords(collectedPopularWords)
        newWeekWords = calculateNewKeywordsForWeek(weekPopularWords, dynamodb)
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

    return &BlogsStatsData{sundayStats, monthStats, newWeekWords}
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

func calculateNewKeywordsForWeek(weekPopularWords string, dynamodb *utils.DynamoDbService) string {
    currentMonthKeywords := dynamodb.GetItem(today.Format(monthFormat))
    currentMonthKeywordsArray := toArray(currentMonthKeywords)
    currentWeekKeywordsArray := toArray(weekPopularWords)
    newWeekWords := findNonDuplicates(currentMonthKeywordsArray, currentWeekKeywordsArray)
    return strings.Join(newWeekWords, ",")
}

func toArray(popularWords string) []string {
    keywords := strings.Split(popularWords, ",")

	for i, keyword := range keywords {
		keywords[i] = strings.TrimSpace(keyword)
	}
    return keywords
}


func findNonDuplicates(arr1 []string, arr2 []string) []string {
	// Create a map to store elements from the first array for efficient lookup
	seen := map[string]bool{}
	for _, element := range arr1 {
		seen[element] = true
	}

	// Iterate through the second array and add non-duplicates to a new slice
	var nonDuplicates []string
	for _, element := range arr2 {
		if !seen[element] {
			nonDuplicates = append(nonDuplicates, element)
		}
	}

	return nonDuplicates
}