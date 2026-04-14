package blogs

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
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
			if dt, ok := dateEl.Attr("datetime"); ok {
				item.Date = parseDate(dt)
			} else {
				item.Date = parseDate(strings.TrimSpace(dateEl.Text()))
			}
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
	time.RFC1123Z,
	time.RFC1123,
	"2006-01-02",
	"January 2, 2006",
	"Jan 2, 2006",
	"02 Jan 2006",
	"2 Jan 2006",
	"January 02, 2006",
	"Mon, 02 Jan 2006",
	"2006/01/02",
}

func parseDate(s string) time.Time {
	s = strings.TrimSpace(s)
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
		`time[datetime]`,
		`time`,
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
	return time.Time{}
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
