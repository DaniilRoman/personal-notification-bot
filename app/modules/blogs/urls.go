package blogs

var blackListKeywords = []string{
	"android", "ios", "redux", "react", "frontend",
	"ui/ux", "career stories", "meeting", "spotlight",
	"internship", "javascript", "css", "html", "typescript",
	"mobile", "uikit", "интерфейс", "дизайн", "мобильн",
	"design", "interface", "A Bootiful Podcast", "Spark", "Docker Desktop",
}

var blogsUrls = []string{
	// Personal blogs
	"https://valeriansaliou.name/blog/rss/",
	"https://jazco.dev/atom.xml",
	"https://vas3k.blog/rss/",
	"https://vladmihalcea.com/blog/feed/",
	"https://piotrminkowski.com/feed/",
	"https://blog.alexellis.io/rss/",
	"https://feeds.feedburner.com/martinkl?format=xml",
	"https://world.hey.com/dhh/feed.atom",
	"https://world.hey.com/jason/feed.atom",
	"https://news.google.com/rss/search?q=site:x.com/levelsio+when:1d&hl=en-US&gl=US&ceid=US:en",
	"https://news.ycombinator.com/rss",
	"https://feeds.feedburner.com/eu-startups",

	// Interesting company blogs
	"https://medium.com/feed/netcracker",
	"https://habr.com/ru/rss/company/just_ai/blog/?fl=ru",
	"https://medium.com/feed/tovieai",
	"https://medium.com/feed/adevinta-tech-blog",

	"https://slack.engineering/feed/",
	"https://engineering.atspotify.com/feed/",
	"https://netflixtechblog.com/feed",
	"https://eng.lyft.com/feed",
// 	"https://stackoverflow.blog/feed",
	"https://medium.com/feed/tinder",
	"https://medium.com/feed/bbc-product-technology",
	"https://open.nytimes.com/feed",
	"https://developers.facebook.com/blog/feed",

	// E-commerce
	"https://engineering.zalando.com/atom.xml",
	"https://habr.com/ru/rss/company/ozontech/blog/?fl=ru",
	"https://habr.com/ru/rss/company/avito/blog/?fl=ru",
	"https://habr.com/ru/rss/company/lamoda/blog/?fl=ru",
	"https://tech.ebayinc.com/rss",
	"https://vinted.engineering/atom.xml",
	"https://www.wix.engineering/blog-feed.xml",
	"https://engineering.squarespace.com/blog?format=rss",
	"https://www.etsy.com/de-en/codeascraft/rss",
	"https://techlab.bol.com/api/v1/en/newsFeed/",
	"https://blog.allegro.tech/feed.xml",

	// Tech products
	"https://github.blog/feed/",
// 	"https://circleci.com/blog/feed.xml",
	"https://kubernetes.io/feed.xml",
	"https://debezium.io/blog.atom",
	"https://www.elastic.co/blog/feed",
	"https://blog.healthchecks.io/feed/",

	// Research labs
	// https://githubnext.com, # didn't find rss feed,
	"http://feeds.feedburner.com/blogspot/gJZg", // Google research
	"https://openai.com/blog/rss.xml",
	"https://research.facebook.com/feed/",
	"https://www.microsoft.com/en-us/research/feed/",
	"https://bair.berkeley.edu/blog/feed",
	"https://www.amazon.science/index.rss",
	"https://www.elastic.co/search-labs/rss/feed",

	// Tech
	"https://aws.amazon.com/blogs/database/tag/dynamodb/feed",
	"https://www.cncf.io/blog/feed/",

	// Others company bogs
	"https://aws.amazon.com/ru/blogs/machine-learning/feed/",
	"https://aws.amazon.com/ru/blogs/architecture/feed/",
	"https://aws.amazon.com/ru/blogs/database/feed/",
	"https://microservices.io/feed.xml",
	"https://www.confluent.io/blog/area/technology/rss.xml",
	"https://blog.cloudflare.com/tag/developers",
	// "https://www.uber.com/en-DE/blog/engineering/rss", // TODO http error: 406 Not Acceptable
	"https://medium.com/feed/miro-engineering",
	"https://habr.com/ru/rss/company/nspk/blog/?fl=ru",
	"https://canvatechblog.com/feed",
	"https://deliveroo.engineering/feed",
	"https://medium.com/feed/paypal-tech",
	"https://medium.com/feed/strava-engineering",
	"https://engineering.linkedin.com/blog.rss.html",
	"https://www.reddit.com/r/RedditEng/.rss",
	"https://doordash.engineering/category/backend/feed/",
	"https://doordash.engineering/category/data-platform/feed/",
	"https://doordash.engineering/category/data-science-and-machine-learning/feed/",

	// Company blogs to delete
	"https://engineering.fb.com/feed/",
	// "https://blog.twitter.com/engineering/en_us/blog.rss", // TODO no published date
	"https://www.datadoghq.com/blog/engineering/index.xml",
	"https://grafana.com/categories/engineering/index.xml",

	// "https://feeds.feedblitz.com/baeldung&x=1",
	// "https://spring.io/blog.atom",
	// "https://www.mongodb.com/blog/rss",
	"https://blog.jetbrains.com/kotlin/category/server/feed/",
	"https://in.relation.to/blog.atom", // Hibernate
	"https://www.hashicorp.com/blog/categories/products-technology/feed.xml",
}
