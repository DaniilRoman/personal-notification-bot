package blogs

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
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
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: &UserAgentTransport{RoundTripper: transport},
		Timeout:   5 * time.Second,
	}
	parser.Client = httpClient

	totalFeeds := len(blogsUrls) + len(scrapeBlogs)
	wg.Add(totalFeeds)
	for _, url := range blogsUrls {
		go parseLastArticle(url, parser, blogsChannel, &wg, client)
	}
	for _, config := range scrapeBlogs {
		go scrapeLastArticle(config, parser, blogsChannel, &wg, client, httpClient)
	}

	go func() {
		wg.Wait()
		close(blogsChannel)
	}()

	titles := []blogUpdate{}
	for title := range blogsChannel {
		titles = append(titles, title)
	}
	log.Println("Finished all: ", titles)

	return &BlogUpdateData{titles}, nil
}

func parseLastArticle(url string, parser *gofeed.Parser, blogs chan<- blogUpdate, wg *sync.WaitGroup, client *utils.ChatGptService) {
	defer wg.Done()
	// log.Println("Starting: " + url)
	// defer log.Println("Finished: " + url)

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

	processFeedItem(lastArticle, client, blogs)
}

func processFeedItem(item *gofeed.Item, client *utils.ChatGptService, blogs chan<- blogUpdate) bool {
	var date time.Time
	if item.PublishedParsed != nil {
		date = *item.PublishedParsed
	}
	if !shouldIncludeArticle(date, item.Title) || isInBlacklist(item) {
		return false
	}
	img, summary, popularWords := getExtraFields(item, client)
	blogs <- NewBlogUpdate(item.Title, item.Link, img, summary, popularWords)
	return true
}

func scrapeLastArticle(config BlogConfig, parser *gofeed.Parser, blogs chan<- blogUpdate, wg *sync.WaitGroup, client *utils.ChatGptService, httpClient *http.Client) {
	defer wg.Done()

	scraper := NewScraper(httpClient)
	rssBytes, items, err := scraper.ScrapeToRSS(config)
	if err != nil {
		log.Printf("Error scraping %s: %s", config.URL, err)
		return
	}
	if len(items) == 0 {
		log.Printf("No items found for %s", config.URL)
		return
	}

	// Save RSS for debugging if environment variable is set
	if os.Getenv("BLOG_DEBUG_RSS") == "1" {
		if err := SaveRSSToFile(rssBytes, config.URL); err != nil {
			log.Printf("Error saving RSS debug file for %s: %s", config.URL, err)
		}
	}

	// Parse the generated RSS with gofeed parser
	feed, err := parser.ParseString(string(rssBytes))
	if err != nil {
		log.Printf("Error parsing generated RSS for %s: %s", config.URL, err)
		return
	}

	if len(feed.Items) == 0 {
		log.Printf("No items in generated RSS for %s", config.URL)
		return
	}

	// Use the first item (RSS items are sorted by date descending)
	processFeedItem(feed.Items[0], client, blogs)
}

func isArticlePublishedYesterday(article *gofeed.Item) bool {
	if article.PublishedParsed == nil {
		return false
	}
	return isArticlePublishedYesterdayTime(*article.PublishedParsed)
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

func isArticlePublishedYesterdayTime(published time.Time) bool {
	return !published.IsZero() && published.After(time.Now().Add(-24*time.Hour))
}

func isTitleInBlacklist(title string) bool {
	return containsInBlacklistKeywords(title)
}

func shouldIncludeArticle(date time.Time, title string) bool {
	if isTitleInBlacklist(title) {
		return false
	}
	// If date is zero, we cannot determine if it's recent, so include it
	if date.IsZero() {
		return true
	}
	return isArticlePublishedYesterdayTime(date)
}

func getExtraFieldsFromLink(link string, client *utils.ChatGptService) (string, string, string) {
	doc, err := utils.GetDoc(link)
	if err != nil {
		log.Printf("Error in getting blog extra data for %s: %s", link, err)
		return "", "", ""
	}
	img := getImage(doc)
	summary, popularWords, err := getSummaryAndSaveStats(doc, client)
	if err != nil {
		log.Printf("Error is during generating a summary for %s: %s", link, err)
	}
	return img, summary, popularWords
}

func getExtraFieldsFromDoc(doc *goquery.Document, client *utils.ChatGptService) (string, string, string) {
	img := getImage(doc)
	summary, popularWords, err := getSummaryAndSaveStats(doc, client)
	if err != nil {
		log.Printf("Error is during generating a summary: %s", err)
	}
	return img, summary, popularWords
}

func getExtraFields(article *gofeed.Item, client *utils.ChatGptService) (string, string, string) {
	return getExtraFieldsFromLink(article.Link, client)
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
		// summary = client.SummarizeText(textToSummarize)
	}()
	go func() {
		defer wg.Done()
		// popularWords = client.ArticlePopularWords(textToSummarize)
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

type UserAgentTransport struct {
	http.RoundTripper
}

func (c *UserAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0")
	return c.RoundTripper.RoundTrip(r)
}
