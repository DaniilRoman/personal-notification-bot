package modules

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
)

var blogsUrls = []string {
	// Personal blogs
	"https://vladmihalcea.com/blog/feed/",
	"https://piotrminkowski.com/feed/",
	"https://blog.alexellis.io/rss/",

	// Interesting company blogs
	"https://feeds.feedblitz.com/baeldung&x=1",
	"https://spring.io/blog.atom",
	"https://blog.jetbrains.com/kotlin/category/server/feed/",
	"https://aws.amazon.com/blogs/aws/feed/",
	"https://medium.com/feed/netcracker",
	"https://habr.com/ru/rss/company/just_ai/blog/?fl=ru",
	"https://medium.com/feed/adevinta-tech-blog",

	"https://slack.engineering/feed/",
	"https://engineering.atspotify.com/feed/",
	"https://engineering.zalando.com/atom.xml",
	"https://github.blog/feed/",
	"https://netflixtechblog.com/feed",
	"https://eng.lyft.com/feed",
	"https://stackoverflow.blog/feed",
	"https://medium.com/feed/tinder",
	"https://medium.com/feed/bbc-product-technology",
	"https://open.nytimes.com/feed",
	"https://neo4j.com/developer-blog/feed/",

	// Research labs
	// https://githubnext.com, # didn't find rss feed,
	"http://feeds.feedburner.com/blogspot/gJZg",  // Google research
	"https://openai.com/blog/rss.xml",
	"https://research.facebook.com/feed/",
	"https://www.microsoft.com/en-us/research/feed/",
	"https://bair.berkeley.edu/blog/feed",

	// Tech
	"https://circleci.com/blog/feed.xml",
	"https://kubernetes.io/feed.xml",
	"https://www.docker.com/blog/feed",
	"https://redis.com/blog/rss",
	"https://www.mongodb.com/blog/rss",
	"https://debezium.io/blog.atom",
	"https://www.elastic.co/blog/feed",
	"https://aws.amazon.com/blogs/database/tag/dynamodb/feed",
	"https://about.gitlab.com/atom.xml",
	"https://feeds.feedburner.com/ContinuousBlog/",  // Jenkins
	"https://in.relation.to/blog.atom",  // Hibernate
	"https://www.cncf.io/blog/feed/",

	// Others company bogs
	"https://www.hashicorp.com/blog/feed.xml",
	"https://microservices.io/feed.xml",
	"https://www.confluent.io/rss.xml",
	"https://blog.cloudflare.com/rss",
	// "https://www.uber.com/en-DE/blog/engineering/rss", // TODO http error: 406 Not Acceptable
	"https://medium.com/feed/miro-engineering",
	"https://habr.com/ru/rss/company/ozontech/blog/?fl=ru",
	"https://habr.com/ru/rss/company/avito/blog/?fl=ru",
	"https://habr.com/ru/rss/company/lamoda/blog/?fl=ru",
	"https://habr.com/ru/rss/company/nspk/blog/?fl=ru",
	"https://canvatechblog.com/feed",
	"https://deliveroo.engineering/feed",
	"https://tech.ebayinc.com/rss",
	"https://medium.com/feed/paypal-tech",
	"https://medium.com/feed/strava-engineering",
	"https://engineering.linkedin.com/blog.rss.html",
	"https://www.reddit.com/r/RedditEng/.rss",

	// Company blogs to delete
	"https://engineering.fb.com/feed/",
	// "https://blog.twitter.com/engineering/en_us/blog.rss", // TODO no published date
	"https://www.datadoghq.com/blog/index.xml",
	"https://grafana.com/categories/engineering/index.xml",
}


func BlogUpdates() *BlogUpdateData {
	blogUpdateData, err := blogUpdates()
    if err != nil {
       log.Printf("Error in blogs updates module: %s", err)
    }
    return blogUpdateData
}

func blogUpdates() (*BlogUpdateData, error) {
	var wg sync.WaitGroup
	titleChannel := make(chan string)
	parser := gofeed.NewParser()

	wg.Add(len(blogsUrls))
	for _, url := range blogsUrls {
		go parseLastArticle(url, parser, titleChannel, &wg)
	}

	go func() {
		wg.Wait()
		close(titleChannel)
	}()

	titles := []string {}
	for title := range titleChannel {
		titles = append(titles, title)
	}
	log.Println("Finished all: ", titles)
	
	return &BlogUpdateData{}, nil
}

func parseLastArticle(url string, parser *gofeed.Parser, job chan<- string, wg *sync.WaitGroup) {
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
		// blogUpdateData := BlogUpdateData()
		// if notInBlacklist(lastArticle) {
		// 	setExtraFields(lastArticle)
		// } else {
		// 	log.Printf("Filtered by topic: %s", lastArticle.Title)
		// }

		// log.Println(lastArticle.Title)
		// log.Println(lastArticle.Link)
		// log.Println("================")

		job <- lastArticle.Title
	}
}

func isArticlePublishedYesterday(article *gofeed.Item) bool {
	return article.PublishedParsed.After(time.Now().Add(-24*time.Hour))
}

func notInBlacklist(article *gofeed.Item) bool {
	return true
}

func setExtraFields(article *gofeed.Item) {

}

type BlogUpdateData struct {

}