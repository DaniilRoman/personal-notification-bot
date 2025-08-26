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
		return err
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
	// Create a context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Give Chrome some time
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var buf []byte
	url := fmt.Sprintf("https://dashboard.simpleanalytics.com/%s.com", projectName)

	// Run tasks
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.ByQuery), // wait until body loads
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
