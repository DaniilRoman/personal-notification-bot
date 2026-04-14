package blogs

var scrapeBlogs = []BlogConfig{
	// Add your non-RSS blog configurations here
	// To debug generated RSS feeds, set environment variable BLOG_DEBUG_RSS=1
	// Generated RSS files will be saved to ./debug_rss/[domain].xml
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
	// IBM Quantum blog (date includes author name, may not parse)
	{
		URL:      "https://www.ibm.com/quantum/blog",
		BaseURL:  "https://www.ibm.com",
		MaxItems: 5,
		Selectors: Selectors{
			Item:  `a[href^="/quantum/blog/"]`,
			Title: "h2, h5",
			Link:  `a[href^="/quantum/blog/"]`,
			Date:  "p.text-caption",
		},
	},
	// Agility Robotics blog
	// {
	// 	URL:      "https://www.agilityrobotics.com/resources?tab=blogs",
	// 	BaseURL:  "https://www.agilityrobotics.com",
	// 	MaxItems: 5,
	// 	Selectors: Selectors{
	// 		Item:  `div[data-w-tab="Blogs"] .collection-item-5.w-dyn-item`,
	// 		Title: ".blog-tease-title",
	// 		Link:  ".blog-tile",
	// 		Date:  ".frame-1307898284 p.blog-tease-meta:nth-child(2)",
	// 	},
	// },
	
	// Sanctuary AI news
	// {
	// 	URL:      "https://www.sanctuary.ai/news/",
	// 	BaseURL:  "https://www.sanctuary.ai",
	// 	MaxItems: 5,
	// 	Selectors: Selectors{
	// 		Item:        "div.summary-block-wrapper:nth-of-type(2) .summary-item",
	// 		Title:       ".summary-title-link",
	// 		Link:        ".summary-title-link",
	// 		Date:        ".summary-metadata-item--date",
	// 		Description: ".summary-excerpt",
	// 	},
	// },
}
