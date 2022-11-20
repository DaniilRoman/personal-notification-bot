from datetime import datetime, timedelta
from time import mktime

import logging
from typing import List

import feedparser

__blacklist_labels = ["android", "ios", "redux", "react", "frontend", "ui/ux", "career stories", "meeting", "spotlight",
                      "internship", "javascript", "css", "html", "typescript", "mobile", "uikit", "интерфейс", "дизайн", 
                      "мобильн", "design", "interface"]

__feed_list = [
    # Personal blogs
    "https://vladmihalcea.com/blog/feed/",
    "https://piotrminkowski.com/feed/",
    "https://blog.alexellis.io/rss/",

    # Interesting company blogs
    "https://spring.io/blog.atom",
    "https://blog.jetbrains.com/kotlin/category/server/feed/",
    "https://aws.amazon.com/blogs/aws/feed/",
    "https://medium.com/feed/netcracker",
    "https://habr.com/ru/rss/company/just_ai/blog/?fl=ru",

    "https://slack.engineering/feed/",
    "https://engineering.atspotify.com/feed/",
    "https://engineering.zalando.com/atom.xml",
    "https://github.blog/category/engineering/feed/",
    "https://netflixtechblog.com/feed",
    "https://eng.lyft.com/feed",
    "https://stackoverflow.blog/engineering/feed/",
    "https://medium.com/feed/tinder",
    "https://medium.com/feed/bbc-product-technology",
    "https://open.nytimes.com/feed",

    # Others company bogs
    "https://www.uber.com/en-DE/blog/engineering/rss",
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
    "https://www.reddit.com/r/RedditEng/.rss"

    # Company blogs to delete
    "https://engineering.fb.com/feed/",
    "https://blog.twitter.com/engineering/en_us/blog.rss",
    "https://www.datadoghq.com/blog/index.xml",
    "https://grafana.com/categories/engineering/index.xml"
]


def __none_blacklist_labels(parsed_article):
    if any(black_list_label in parsed_article.title.lower() for black_list_label in __blacklist_labels):
        return False
    if parsed_article.get("tags") is not None:
        for tag in parsed_article.tags:
            if any(black_list_label == tag.term.lower() for black_list_label in __blacklist_labels):
                return False

    return True


def __add_new_article_to_res_list(feed, res: List[str]):
    logging.info(f"Start process {feed}")
    parsed = feedparser.parse(feed)
    logging.info(f"Parsed {feed}")

    if len(parsed.entries) == 0:
        logging.error(f"Cannot parce {feed}")
        return

    last_article = parsed.entries[0]

    if last_article.get("published_parsed") is None:
        article_published = datetime.fromtimestamp(mktime(parsed.updated_parsed)).date()
    else:
        article_published = datetime.fromtimestamp(mktime(last_article.published_parsed)).date()
    prev_day = datetime.today().date() - timedelta(days=1)

    if article_published == prev_day:
        res_article_str = f"{last_article.title} {last_article.link}"
        if __none_blacklist_labels(last_article):
            res.append(res_article_str)
        else:
            logging.info(f"Filtered by topic: {res_article_str}")


def _blog_updates():
    res = []
    for feed in __feed_list:
        try:
            __add_new_article_to_res_list(feed, res)
        except:
            logging.exception(f"Couldn't process blog: {str(feed)}")

    return "\n".join(res)


def blog_updates():
    try:
        return _blog_updates()
    except:
        logging.exception("Couldn't get blog posts")
        return "Couldn't get blog posts"
