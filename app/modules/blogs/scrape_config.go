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
	// Figure AI news (dates in datetime attribute)
	{
		URL:      "https://www.figure.ai/news",
		BaseURL:  "https://www.figure.ai",
		MaxItems: 5,
		Selectors: Selectors{
			Item:        "a.article-list-item",
			Title:       "h1.article-list-item__heading",
			Link:        "a.article-list-item",
			Date:        "time.article-list-item__publication-date",
			Description: "",
		},
	},
	// NEURA Robotics news (dates in time element)
	{
		URL:      "https://neura-robotics.com/news/",
		BaseURL:  "https://neura-robotics.com",
		MaxItems: 5,
		Selectors: Selectors{
			Item:        ".swiper-slide.e-loop-item",
			Title:       "h2.elementor-heading-title a",
			Link:        "h2.elementor-heading-title a",
			Date:        "time",
			Description: "",
		},
	},
	// HRL Uni Bonn news (dates not available)
	{
		URL:      "https://www.hrl.uni-bonn.de/api/news",
		BaseURL:  "https://www.hrl.uni-bonn.de",
		MaxItems: 5,
		Selectors: Selectors{
			Item:        ".portletNavigationTree .navTree li.navTreeItem",
			Title:       "a.state-published",
			Link:        "a.state-published",
			Date:        "", // Dates not available on this site
			Description: "",
		},
	},
	// Agility Robotics blog
	{
		URL:      "https://www.agilityrobotics.com/resources?tab=blogs",
		BaseURL:  "https://www.agilityrobotics.com",
		MaxItems: 5,
		Selectors: Selectors{
			Item:  `div[data-w-tab="Blogs"] .collection-item-5.w-dyn-item`,
			Title: ".blog-tease-title",
			Link:  ".blog-tile",
			Date:  ".blog-tile-copy-wrapper .frame-1307898284 p.blog-tease-meta",
		},
	},
	// Sanctuary AI news
	{
		URL:      "https://www.sanctuary.ai/news/",
		BaseURL:  "https://www.sanctuary.ai",
		MaxItems: 5,
		Selectors: Selectors{
			Item:        "div.summary-block-wrapper:nth-of-type(2) .summary-item",
			Title:       ".summary-title-link",
			Link:        ".summary-title-link",
			Date:        ".summary-metadata.summary-metadata--primary time.summary-metadata-item.summary-metadata-item--date",
			Description: ".summary-excerpt",
		},
	},
}
