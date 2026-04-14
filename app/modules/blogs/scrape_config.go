package blogs

var scrapeBlogs = []BlogConfig{
	// Example configuration:
	// {
	// 	URL:      "https://example.com/blog",
	// 	BaseURL:  "https://example.com",
	// 	MaxItems: 10,
	// 	Selectors: Selectors{
	// 		Item:        "article",
	// 		Title:       "h2 a",
	// 		Link:        "h2 a",
	// 		Date:        "time",
	// 		Description: ".excerpt",
	// 	},
	// },
	// Boston Dynamics blog (no RSS feed)
	{
		URL:      "https://bostondynamics.com/blog/",
		BaseURL:  "https://bostondynamics.com",
		MaxItems: 5,
		Selectors: Selectors{
			Item:  "article.PostAjaxFilter-card",
			Title: "p.PostAjaxFilter-card-title a",
			Link:  "p.PostAjaxFilter-card-title a",
			Date:  "", // Date extracted from article page
		},
	},
	// Add your non-RSS blog configurations here
	// To debug generated RSS feeds, set environment variable BLOG_DEBUG_RSS=1
	// Generated RSS files will be saved to ./debug_rss/[domain].xml
}
