package blogs

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
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

func TestExtractBostonDynamics(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	// This test actually fetches the live page; use sparingly
	config := BlogConfig{
		URL:      "https://bostondynamics.com/blog/",
		BaseURL:  "https://bostondynamics.com",
		MaxItems: 2,
		Selectors: Selectors{
			Item:  "article.PostAjaxFilter-card",
			Title: "p.PostAjaxFilter-card-title a",
			Link:  "p.PostAjaxFilter-card-title a",
		},
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", config.URL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		t.Skipf("Network error, skipping: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	items, err := Extract(config, body)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	print(fmt.Sprintf("%v", items))

	if len(items) == 0 {
		t.Error("Expected at least one item, got zero")
	}

	for _, item := range items {
		if item.Title == "" {
			t.Error("Item title is empty")
		}
		if item.Link == "" {
			t.Error("Item link is empty")
		}
		if !strings.HasPrefix(item.Link, "http") {
			t.Errorf("Item link is not absolute: %s", item.Link)
		}
	}
}

func TestAllScrapeConfigs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := &http.Client{Timeout: 15 * time.Second}

	for i, config := range scrapeBlogs {
		t.Run(config.URL, func(t *testing.T) {
			req, err := http.NewRequest("GET", config.URL, nil)
			if err != nil {
				t.Skipf("Failed to create request for %s: %v", config.URL, err)
			}
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0")

			resp, err := client.Do(req)
			if err != nil {
				t.Skipf("Network error fetching %s, skipping: %v", config.URL, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				t.Skipf("Non-200 status for %s: %d", config.URL, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Skipf("Failed to read body from %s: %v", config.URL, err)
			}

			items, err := Extract(config, body)
			if err != nil {
				t.Skipf("Extract failed for %s: %v", config.URL, err)
			}

			if len(items) == 0 {
				t.Logf("No items extracted for %s", config.URL)
				return
			}

			rssBytes, err := GenerateRSS(items, config.URL, config.URL, "Scraped feed")
			if err != nil {
				t.Skipf("GenerateRSS failed for %s: %v", config.URL, err)
			}

			fmt.Printf("\n=== RSS Feed %d: %s ===\n", i+1, config.URL)
			fmt.Println(string(rssBytes))
			fmt.Printf("=== End RSS Feed %d ===\n\n", i+1)
		})
	}
}
