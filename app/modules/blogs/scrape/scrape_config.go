package scrape

var ScrapeBlogs = []BlogConfig{
	// Add your non-RSS blog configurations here
	// To debug generated RSS feeds, set environment variable BLOG_DEBUG_RSS=1
	// Generated RSS files will be saved to ./debug_rss/[domain].xml
	// Example configuration:
	// {
	// 	URL:      "https://example.com/blog",
	// 	BaseURL:  "https://example.com",
	// 	MaxItems: 3,
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
		MaxItems: 3,
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
		MaxItems: 3,
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
		MaxItems: 3,
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
		MaxItems: 3,
		Selectors: Selectors{
			Item:        ".swiper-slide.e-loop-item",
			Title:       "h2.elementor-heading-title a",
			Link:        "h2.elementor-heading-title a",
			Date:        "time",
			Description: "",
		},
	},
	// HRL Uni Bonn news (dates not available) - temporarily disabled due to Accept header issue
	// {
	// 	URL:      "https://www.hrl.uni-bonn.de/api/news",
	// 	BaseURL:  "https://www.hrl.uni-bonn.de",
	// 	MaxItems: 3,
	// 	Selectors: Selectors{
	// 		Item:        ".portletNavigationTree .navTree li.navTreeItem",
	// 		Title:       "a.state-published",
	// 		Link:        "a.state-published",
	// 		Date:        "", // Dates not available on this site
	// 		Description: "",
	// 	},
	// },
	// Agility Robotics blog
	{
		URL:      "https://www.agilityrobotics.com/resources?tab=blogs",
		BaseURL:  "https://www.agilityrobotics.com",
		MaxItems: 3,
		Selectors: Selectors{
			Item:  `div[data-w-tab="Blogs"] .collection-item-5.w-dyn-item`,
			Title: ".blog-tease-title",
			Link:  ".blog-tile",
			Date:  ".blog-tile-copy-wrapper .frame-1307898284 p.blog-tease-meta:nth-of-type(2)",
		},
	},
	// Sanctuary AI news
	{
		URL:      "https://www.sanctuary.ai/news/",
		BaseURL:  "https://www.sanctuary.ai",
		MaxItems: 3,
		Selectors: Selectors{
			Item:        ".summary-item",
			Title:       ".summary-title-link",
			Link:        ".summary-title-link",
			Date:        ".summary-metadata.summary-metadata--primary time.summary-metadata-item.summary-metadata-item--date",
			Description: ".summary-excerpt",
		},
	},
	// Agile Robots news
	{
		URL:      "https://www.agile-robots.com/en/news/",
		BaseURL:  "https://www.agile-robots.com",
		MaxItems: 3,
		Selectors: Selectors{
			Item:        ".article-list-teaser-element",
			Title:       "h3.title",
			Link:        "a[href]",
			Date:        "time[datetime]",
			Description: "",
		},
	},
	// DFKI Robotics Innovation Center news
	{
		URL:      "https://robotik.dfki-bremen.de/en/startpage",
		BaseURL:  "https://robotik.dfki-bremen.de",
		MaxItems: 3,
		Selectors: Selectors{
			Item:        ".news-teaser-item",
			Title:       "h4.news-teaser-title a",
			Link:        "h4.news-teaser-title a[href]",
			Date:        "time[datetime]",
			Description: "",
		},
	},
	// TUM MIRMI news
	{
		URL:      "https://www.mirmi.tum.de/en/mirmi/news/",
		BaseURL:  "https://www.mirmi.tum.de",
		MaxItems: 3,
		Selectors: Selectors{
			Item:        ".article.articletype-0",
			Title:       "h3 a",
			Link:        "h3 a[href]",
			Date:        "time[datetime]",
			Description: ".teaser-text p",
		},
	},
	// Alice & Bob blog
	{
		URL:      "https://alice-bob.com/blog/",
		BaseURL:  "https://alice-bob.com",
		MaxItems: 3,
		Selectors: Selectors{
			Item:        "li.c-archive__article",
			Title:       "h4.c-post-card__title",
			Link:        "a.c-post-card__inner[href]",
			Date:        ".c-post-card__meta-item",
			Description: "",
		},
	},
	// Quantinuum blog
	{
		URL:      "https://www.quantinuum.com/news/blog#",
		BaseURL:  "https://www.quantinuum.com",
		MaxItems: 3,
		Selectors: Selectors{
			Item:        ".blog_cms_item",
			Title:       "div[fs-cmsfilter-field='heading']",
			Link:        "a.textlink[href]",
			Date:        ".eyebrow_content_contain > .blog_eyebrow",
			Description: ".blog-content-filter p:first-child",
		},
	},
	// Pasqal blog
	{
		URL:      "https://www.pasqal.com/blog/",
		BaseURL:  "https://www.pasqal.com",
		MaxItems: 3,
		Selectors: Selectors{
			Item:        "div.featured-post",
			Title:       "h2.post-title",
			Link:        "a.post-card",
			Date:        "time.entry-date",
			Description: "",
		},
	},
	// Waymo blog (latest posts)
	{
		URL:      "https://waymo.com/blog/",
		BaseURL:  "https://waymo.com",
		MaxItems: 3,
		Selectors: Selectors{
			Item:        "li._postItem_1n64j_107",
			Title:       "h2._postTitle_1n64j_120",
			Link:        "a._postLink_1n64j_400",
			Date:        "p._postDate_1n64j_378",
			Description: "div._postSummary_1n64j_133",
		},
	},
	// Covariant AI blog
	{
		URL:      "https://covariant.ai/resources/blog/",
		BaseURL:  "https://covariant.ai",
		MaxItems: 3,
		Selectors: Selectors{
			Item:        `a.sanity-cta.internal.custom:not(.mini):has(figure.sanity-img)`,
			Title:       `.info .upper .title`,
			Link:        `a.sanity-cta.internal.custom:not(.mini):has(figure.sanity-img)`,
			Date:        "", // Date not present in listing
			Description: `.info .upper .excerpt`,
		},
	},
	// Starship Technologies newsroom
	{
		URL:      "https://www.starship.xyz/newsroom/",
		BaseURL:  "https://www.starship.xyz",
		MaxItems: 3,
		Selectors: Selectors{
			Item:        ".news-card",
			Title:       ".news-card__title a",
			Link:        ".news-card__title a",
			Date:        ".news-card__date",
			Description: "",
		},
	},
}
