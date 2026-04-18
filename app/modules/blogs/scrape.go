package blogs

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"main/utils"
)

type BlogConfig struct {
	URL       string
	BaseURL   string
	Selectors Selectors
	MaxItems  int
}

type Selectors struct {
	Item        string
	Title       string
	Link        string
	Date        string
	Description string
}

type ScrapedItem struct {
	Title       string
	Link        string
	Date        time.Time
	Description string
}

type Scraper struct {
	Client    *http.Client
	UserAgent string
}

func NewScraper(client *http.Client) *Scraper {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	return &Scraper{
		Client:    client,
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0",
	}
}

func (s *Scraper) fetchDoc(url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", s.UserAgent)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

func (s *Scraper) ScrapeToRSS(config BlogConfig) ([]byte, []ScrapedItem, error) {
	doc, err := s.fetchDoc(config.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("fetching %s: %w", config.URL, err)
	}

	// Get HTML body bytes for Extract
	html, err := doc.Html()
	if err != nil {
		return nil, nil, fmt.Errorf("getting HTML: %w", err)
	}

	items, err := Extract(config, []byte(html))
	if err != nil {
		return nil, nil, fmt.Errorf("extracting items: %w", err)
	}

	// Fetch missing dates using scraper's fetchDoc
	items = fetchMissingDatesWithFetcher(items, s.fetchDoc)

	rssBytes, err := GenerateRSS(items, config.URL, config.URL, "Scraped feed")
	if err != nil {
		return nil, nil, fmt.Errorf("generating RSS: %w", err)
	}

	return rssBytes, items, nil
}

func Extract(config BlogConfig, body []byte) ([]ScrapedItem, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	sel := config.Selectors
	var items []ScrapedItem

	doc.Find(sel.Item).EachWithBreak(func(i int, s *goquery.Selection) bool {
		if i >= config.MaxItems {
			return false
		}

		item := ScrapedItem{}

		if sel.Title != "" {
			item.Title = strings.TrimSpace(s.Find(sel.Title).First().Text())
		}
		if item.Title == "" {
			item.Title = strings.TrimSpace(s.Text())
		}

		if sel.Link != "" {
			linkEl := s.Find(sel.Link).First()
			href, exists := linkEl.Attr("href")
			if !exists {
				href, _ = s.Attr("href")
			}
			item.Link = resolveURL(config.BaseURL, config.URL, href)
		}

		if sel.Date != "" {
			dateEl := s.Find(sel.Date).First()
			var rawDate string
			if dt, ok := dateEl.Attr("datetime"); ok {
				rawDate = dt
				item.Date = parseDate(dt)
			} else {
				rawDate = strings.TrimSpace(dateEl.Text())
				item.Date = parseDate(rawDate)
			}
			if item.Date.IsZero() && rawDate != "" {
				log.Printf("Date selector '%s' found raw date '%s' but failed to parse", sel.Date, rawDate)
			}
		} else {
			log.Printf("No date selector for config %s", config.URL)
		}

		if sel.Description != "" {
			item.Description = strings.TrimSpace(s.Find(sel.Description).First().Text())
		}

		items = append(items, item)
		return true
	})

	return items, nil
}

func resolveURL(baseURL, pageURL, href string) string {
	href = strings.TrimSpace(href)
	if href == "" {
		return ""
	}
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}
	root := baseURL
	if root == "" {
		root = pageURL
	}
	base, err := url.Parse(root)
	if err != nil {
		return href
	}
	ref, err := url.Parse(href)
	if err != nil {
		return href
	}
	return base.ResolveReference(ref).String()
}

var dateFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	time.RFC1123Z,
	time.RFC1123,
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05-07:00",
	"2006-01-02T15:04:05-0700",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02",
	"2006.01.02",
	"02.01.2006",
	"2.1.2006",
	"02.01.06",
	"2.1.06",
	"January 2, 2006",
	"Jan 2, 2006",
	"Jan 2, 2006 3:04 PM",
	"Jan 2, 2006 15:04",
	"02 Jan 2006",
	"2 Jan 2006",
	"January 02, 2006",
	"Mon, 02 Jan 2006",
	"Monday, January 2, 2006",
	"Mon Jan 2 2006",
	"2006/01/02",
	"2006/01/02 15:04:05",
}

func cleanDateString(s string) string {
	original := s
	s = strings.TrimSpace(s)
	// Remove common prefixes (longer first)
	prefixes := []string{"Published on", "Posted on", "Published", "Posted", "Updated", "Date:"}
	for _, p := range prefixes {
		if strings.HasPrefix(strings.ToLower(s), strings.ToLower(p)) {
			// Remove prefix using slice (case-insensitive)
			s = s[len(p):]
			s = strings.TrimSpace(s)
			// Also remove optional colon
			s = strings.TrimPrefix(s, ":")
			s = strings.TrimSpace(s)
			break // only remove one prefix
		}
	}
	// Split on common separators that indicate extra metadata
	separators := []string{" • ", " | ", " - ", " – ", " — "}
	for _, sep := range separators {
		if idx := strings.Index(s, sep); idx > 0 {
			s = s[:idx]
			break
		}
	}
	// Remove extra spaces and newlines
	s = strings.Join(strings.Fields(s), " ")
	if os.Getenv("BLOG_DEBUG_DATE") == "1" {
		log.Printf("cleanDateString: input=%q, output=%q", original, s)
	}
	return s
}

func parseDate(s string) time.Time {
	s = cleanDateString(s)
	if s == "" {
		return time.Time{}
	}
	for _, f := range dateFormats {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

func ExtractDateFromDoc(doc *goquery.Document) time.Time {
	// Try meta tags
	selectors := []string{
		`meta[property="article:published_time"]`,
		`meta[name="article:published_time"]`,
		`meta[property="og:published_time"]`,
		`meta[name="datePublished"]`,
		`meta[property="datePublished"]`,
		`meta[property="og:updated_time"]`,
		`meta[name="publish_date"]`,
		`time[datetime]`,
		`time`,
		`span.sanity-date`,
		`[class*="date"]`,
	}
	for _, sel := range selectors {
		selection := doc.Find(sel).First()
		if selection.Length() == 0 {
			continue
		}
		if sel == `time[datetime]` || strings.Contains(sel, "meta") {
			if dt, ok := selection.Attr("datetime"); ok {
				if t := parseDate(dt); !t.IsZero() {
					return t
				}
			}
			if dt, ok := selection.Attr("content"); ok {
				if t := parseDate(dt); !t.IsZero() {
					return t
				}
			}
		}
		// For plain time element, use text
		if t := parseDate(selection.Text()); !t.IsZero() {
			return t
		}
	}
	// Try JSON-LD script tags
	var found time.Time
	doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
		if !found.IsZero() {
			return
		}
		text := s.Text()
		// Look for datePublished field
		re := regexp.MustCompile(`"datePublished"\s*:\s*"([^"]+)"`)
		if matches := re.FindStringSubmatch(text); matches != nil {
			if t := parseDate(matches[1]); !t.IsZero() {
				found = t
			}
		}
	})
	if !found.IsZero() {
		return found
	}
	return time.Time{}
}

func fetchMissingDatesWithFetcher(items []ScrapedItem, fetchDoc func(string) (*goquery.Document, error)) []ScrapedItem {
	updated := make([]ScrapedItem, len(items))
	for i, item := range items {
		updated[i] = item
		if !item.Date.IsZero() {
			continue
		}
		if item.Link == "" {
			continue
		}
		log.Printf("Fetching date from article page: %s", item.Link)
		doc, err := fetchDoc(item.Link)
		if err != nil {
			log.Printf("Failed to fetch article page %s: %v", item.Link, err)
			continue
		}
		date := ExtractDateFromDoc(doc)
		if date.IsZero() {
			log.Printf("No date found in article %s", item.Link)
			continue
		}
		updated[i].Date = date
		log.Printf("Found date %v for %s", date, item.Title)
	}
	return updated
}

func fetchMissingDates(items []ScrapedItem) []ScrapedItem {
	return fetchMissingDatesWithFetcher(items, utils.GetDoc)
}

func GenerateRSS(items []ScrapedItem, feedTitle, feedLink, feedDescription string) ([]byte, error) {
	// Sort items by date descending (most recent first)
	sort.Slice(items, func(i, j int) bool {
		if items[i].Date.IsZero() && items[j].Date.IsZero() {
			return false // keep order
		}
		if items[i].Date.IsZero() {
			return false // zero dates go to end
		}
		if items[j].Date.IsZero() {
			return true // non-zero before zero
		}
		return items[i].Date.After(items[j].Date)
	})

	now := time.Now()
	feed := &feeds.Feed{
		Title:       feedTitle,
		Link:        &feeds.Link{Href: feedLink},
		Description: feedDescription,
		Created:     now,
	}

	var feedItems []*feeds.Item
	for _, item := range items {
		feedItem := &feeds.Item{
			Title:       item.Title,
			Link:        &feeds.Link{Href: item.Link},
			Description: item.Description,
		}
		if !item.Date.IsZero() {
			feedItem.Created = item.Date
		}
		// If date is zero, leave Created as zero time (will be omitted from RSS)
		feedItems = append(feedItems, feedItem)
	}
	feed.Items = feedItems

	rssString, err := feed.ToRss()
	if err != nil {
		return nil, err
	}
	return []byte(rssString), nil
}

func SaveRSSToFile(rssBytes []byte, urlStr string) error {
	debugDir := "./debug_rss"
	if err := os.MkdirAll(debugDir, 0755); err != nil {
		return err
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	filename := strings.ReplaceAll(u.Hostname(), ".", "_") + ".xml"
	filepath := filepath.Join(debugDir, filename)
	return os.WriteFile(filepath, rssBytes, 0644)
}
