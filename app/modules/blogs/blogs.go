package blogs

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"main/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"jaytaylor.com/html2text"
)


func BlogUpdates(client *utils.ChatGptService) *BlogUpdateData {
	blogUpdateData, err := blogUpdates(client)
    if err != nil {
       log.Printf("Error in blogs updates module: %s", err)
    }
    return blogUpdateData
}

func blogUpdates(client *utils.ChatGptService) (*BlogUpdateData, error) {
	var wg sync.WaitGroup
	blogsChannel := make(chan blogUpdate)
	parser := gofeed.NewParser()

	wg.Add(len(blogsUrls))
	for _, url := range blogsUrls {
		go parseLastArticle(url, parser, blogsChannel, &wg, client)
	}

	go func() {
		wg.Wait()
		close(blogsChannel)
	}()

	titles := []blogUpdate {}
	for title := range blogsChannel {
		titles = append(titles, title)
	}
	log.Println("Finished all: ", titles)
	
	return &BlogUpdateData{titles}, nil
}

func parseLastArticle(url string, parser *gofeed.Parser, blogs chan<- blogUpdate, wg *sync.WaitGroup, client *utils.ChatGptService) {
	defer wg.Done()
	log.Println("Starting: " + url)
	defer log.Println("Finished: " + url)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	feed, err := parser.ParseURLWithContext(url, ctx)
	if err != nil {
		log.Printf("Error happednd: %s", err)
		return
	}

	if len(feed.Items) == 0 {
		log.Printf("No items: %s", url)
		return
	}
	lastArticle := feed.Items[0]

	if !isArticlePublishedYesterday(lastArticle) || isInBlacklist(lastArticle) {
		log.Printf("Article %s is filtered out", lastArticle.Title)
		return
	}

	img, summary, popularWords := getExtraFields(lastArticle, client)

	blogs <- NewBlogUpdate(lastArticle.Title, lastArticle.Link, img, summary, popularWords)
}

func isArticlePublishedYesterday(article *gofeed.Item) bool {
	return article.PublishedParsed.After(time.Now().Add(-24*time.Hour))
}

func isInBlacklist(article *gofeed.Item) bool {
	if containsInBlacklistKeywords(article.Title) {
		return true
	}

	if article.Categories != nil {
		for _, tag := range article.Categories {
			if containsInBlacklistKeywords(tag) {
				return true
			}
		}
	}

	return false
}

func containsInBlacklistKeywords(s string) bool {
	for _, blacklistLabel := range blackListKeywords {
		if strings.Contains(strings.ToLower(s), blacklistLabel) {
			return true
		}
	}
	return false
}

func getExtraFields(article *gofeed.Item, client *utils.ChatGptService) (string, string, string) {
	doc, err := utils.GetDoc(article.Link)
	if err != nil {
		log.Printf("Error in getting blog extra data for %s: %s", article.Link, err)
		return "", "", ""
	}
	img := getImage(doc)
	summary, popularWords, err := getSummaryAndSaveStats(doc, client)
	if err != nil {
		log.Printf("Error is during generating a summary for %s: %s", article.Link, err)
	}
	return img, summary, popularWords
}

func getImage(doc *goquery.Document) string {
	img, _ := doc.Find(`meta[property="og:image"]`).Attr("content")
	return img
}

func getSummaryAndSaveStats(doc *goquery.Document, client *utils.ChatGptService) (string, string, error) {
	textToSummarize, err := htmlToText(doc)
	textToSummarize = filterLongLines(textToSummarize)
	if err != nil {
		return "", "", err
	}
	textToSummarize = strings.ReplaceAll(textToSummarize, "\n", " ")
	textToSummarize = strings.ReplaceAll(textToSummarize, "\t", " ")
	textToSummarize = strings.ReplaceAll(textToSummarize, "  ", " ")

	var wg sync.WaitGroup
	wg.Add(2)
	
	summary := ""
	popularWords := ""
	go func() {
		defer wg.Done()
		summary = client.SummarizeText(textToSummarize)
	}()
	go func() {
		defer wg.Done()
		popularWords = client.ArticlePopularWords(textToSummarize)
	}()
	wg.Wait()
	
	return summary, popularWords, nil
}

func htmlToText(doc *goquery.Document) (string, error) {
	html, err := doc.Html()
	if err != nil {
		return "", err
	}
	return html2text.FromString(html, html2text.Options{PrettyTables: false, TextOnly: true, OmitLinks: true})
}

func filterLongLines(text string) string {
	var filteredLines []string
	for _, line := range strings.Split(text, "\n") {
		if len(line) > 42 {
			filteredLines = append(filteredLines, line)
		}
	}
	return strings.Join(filteredLines, " ")
}
