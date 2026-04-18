package blogs

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

func TestExtractDateFromDoc(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected time.Time
	}{
		{
			name:     "article published meta",
			html:     `<meta property="article:published_time" content="2026-04-14T17:37:25+00:00">`,
			expected: time.Date(2026, 4, 14, 17, 37, 25, 0, time.UTC),
		},
		{
			name:     "datetime attribute",
			html:     `<time datetime="2026-04-14">April 14, 2026</time>`,
			expected: time.Date(2026, 4, 14, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "no date",
			html:     `<div>No date here</div>`,
			expected: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			if err != nil {
				t.Fatalf("Failed to parse HTML: %v", err)
			}
			got := ExtractDateFromDoc(doc)
			if !got.Equal(tt.expected) {
				t.Errorf("ExtractDateFromDoc() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAllScrapeConfigs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := &http.Client{Timeout: 15 * time.Second}

	for i, config := range scrapeBlogs {
		t.Run(config.URL, func(t *testing.T) {
			scraper := NewScraper(client)
			rssBytes, items, err := scraper.ScrapeToRSS(config)
			if err != nil {
				t.Skipf("ScrapeToRSS failed for %s: %v", config.URL, err)
			}

			if len(items) == 0 {
				t.Errorf("No items extracted for %s", config.URL)
				return
			}

			// Verify all items have title, link, and date
			for _, item := range items {
				if item.Title == "" {
					t.Errorf("Item missing title: %s", item.Link)
				}
				if item.Link == "" {
					t.Errorf("Item missing link: %q", item.Title)
				}
				if !strings.HasPrefix(item.Link, "http") {
					t.Errorf("Item link is not absolute: %q (%s)", item.Title, item.Link)
				}
				if item.Date.IsZero() {
					t.Errorf("Item missing date: %q (%s)", item.Title, item.Link)
				}
			}

			// Validate RSS parses correctly
			parser := gofeed.NewParser()
			feed, err := parser.ParseString(string(rssBytes))
			if err != nil {
				t.Errorf("Generated RSS doesn't parse: %v", err)
			} else if len(feed.Items) != len(items) {
				t.Errorf("RSS item count mismatch: expected %d, got %d", len(items), len(feed.Items))
			}

			fmt.Printf("\n=== RSS Feed %d: %s ===\n", i+1, config.URL)
			fmt.Println(string(rssBytes))
			fmt.Printf("=== End RSS Feed %d ===\n\n", i+1)
		})
	}
}

func TestParseDateCleaning(t *testing.T) {
	tests := []struct {
		input       string
		shouldParse bool
	}{
		{"16 Mar 2026 • Leron Gil,  Julianna Roberts, Sanskriti Deva, Robert Davis", true},
		{"23 Mar 2026 • Rafi Letzter", true},
		{"Published: January 2, 2026", true},
		{"Posted on 2026-01-02", true},
		{"Date: 2026/01/02 15:04:05", true},
		{"2026-01-02T15:04:05Z", true},
		{"", false},
		{"invalid date", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseDate(tt.input)
			if tt.shouldParse && got.IsZero() {
				t.Errorf("parseDate(%q) = zero time, expected non-zero", tt.input)
			}
			if !tt.shouldParse && !got.IsZero() {
				t.Errorf("parseDate(%q) = non-zero time %v, expected zero", tt.input, got)
			}
		})
	}
}
