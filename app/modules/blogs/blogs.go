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
	"github.com/sashabaranov/go-openai"
)


func BlogUpdates(client *openai.Client) *BlogUpdateData {
	blogUpdateData, err := blogUpdates(client)
    if err != nil {
       log.Printf("Error in blogs updates module: %s", err)
    }
    return blogUpdateData
}

func blogUpdates(client *openai.Client) (*BlogUpdateData, error) {
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

func parseLastArticle(url string, parser *gofeed.Parser, blogs chan<- blogUpdate, wg *sync.WaitGroup, client *openai.Client) {
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

	if isArticlePublishedYesterday(lastArticle) {
		if isInBlacklist(lastArticle) {
			log.Printf("Article %s in black list", lastArticle.Title)
			return
		}
		img, summary := getExtraFields(lastArticle, client)

		blogs <- NewBlogUpdate(lastArticle.Title, lastArticle.Link, img, summary)
	}
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

func getExtraFields(article *gofeed.Item, client *openai.Client) (string, string) {
	doc, err := utils.GetDoc(article.Link)
	if err != nil {
		log.Printf("Error in getting blog extra data for %s: %s", article.Link, err)
		return "", ""
	}
	img := getImage(doc)
	summary := getSummary(doc, client)
	return img, summary
}

func getImage(doc *goquery.Document) string {
	img, _ := doc.Find(`meta[property="og:image"]`).Attr("content")
	return img
}

func getSummary(doc *goquery.Document, client *openai.Client) string {
	textToSummarize := doc.Text()
	textToSummarize = strings.ReplaceAll(textToSummarize, "\n", " ")
	textToSummarize = strings.ReplaceAll(textToSummarize, "\t", " ")
	textToSummarize = strings.ReplaceAll(textToSummarize, "  ", " ")

	return utils.SummarizeText(textToSummarize, client)
}

