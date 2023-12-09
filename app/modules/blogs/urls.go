package blogs


var blackListKeywords = []string {
	"android", "ios", "redux", "react", "frontend", 
	"ui/ux", "career stories", "meeting", "spotlight",
    "internship", "javascript", "css", "html", "typescript", 
	"mobile", "uikit", "интерфейс", "дизайн", "мобильн", 
	"design", "interface", "A Bootiful Podcast", "Spark", "Docker Desktop",
}

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
