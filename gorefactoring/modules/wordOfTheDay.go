package modules

import (
	"fmt"
	"log"
	utils "main/utils"
)

func WordOfTheDay() *WordOfTheDayData {
	res, err := wordOfTheDay()
    if err != nil {
		log.Printf("Error in Word of the day module: %s", err)
	 }
	return res
}

func wordOfTheDay() (*WordOfTheDayData, error) {
	wordOfTheDayURL := "https://www.nytimes.com/column/learning-word-of-the-day"
	doc, err := utils.GetDoc(wordOfTheDayURL)
	if err != nil {
		return nil, err
	}

	lastArticle := doc.Find("article")
	lastArticleURL, _ := lastArticle.Find("a").First().Attr("href")
	lastArticleTitle := lastArticle.Find("h3").First().Text()

	res := fmt.Sprintf("[%s](https://www.nytimes.com%s)\n", lastArticleTitle, lastArticleURL)

	return &WordOfTheDayData{res}, nil
}

type WordOfTheDayData struct {
	data string
}

func (d *WordOfTheDayData) String() string {
	if d == nil {
		return ""
	}
    return d.data
}