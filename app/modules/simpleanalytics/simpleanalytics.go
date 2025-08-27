package simpleanalytics

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

func MakeSimpleAnalyticsScreenshot() {
	if err := os.MkdirAll("./data/images", 0755); err != nil {
		log.Printf("Error in SimpleAnalytics while creating the folders: %s", err)
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		makeSimpleAnalyticsScreenshot("daniilroman")
	}()

	go func() {
		defer wg.Done()
		makeSimpleAnalyticsScreenshot("timelinemap.daniilroman")
	}()

	go func() {
		defer wg.Done()
		makeSimpleAnalyticsScreenshot("substacktrends.daniilroman")
	}()

	wg.Wait()
}

func makeSimpleAnalyticsScreenshot(projectName string) {
	// Configure options for chrome
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// Disable sandbox in CI environment
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		// Add additional options for stability
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("headless", true),
	)

	// Create allocator context
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create browser context
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout for the entire operation
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var buf []byte
	url := fmt.Sprintf("https://dashboard.simpleanalytics.com/%s.com", projectName)

	// Run tasks
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`[class*="graph-container"]`, chromedp.ByQuery),
		chromedp.FullScreenshot(&buf, 90),
	); err != nil {
		log.Fatal(err)
	}

	// Save screenshot
	if err := os.WriteFile(fmt.Sprintf("./data/images/%s.png", projectName), buf, 0644); err != nil {
		log.Fatal(err)
	}

	log.Println(fmt.Sprintf("Screenshot saved as %s.png", projectName))
}
