from datetime import datetime, timedelta
from time import mktime

import logging
import feedparser


def _blog_updates():
    feed_list = [
        "https://vladmihalcea.com/blog/feed/",
        "https://piotrminkowski.com/feed/",
        "https://blog.alexellis.io/rss/",
        "https://slack.engineering/feed/",
        "https://engineering.fb.com/feed/",
        "https://engineering.atspotify.com/feed/",
        "https://engineering.zalando.com/atom.xml",
#         "https://blog.twitter.com/engineering/en_us/blog.rss", # article_published = datetime.fromtimestamp(mktime(last_article.published_parsed)).date() // TypeError: Tuple or struct_time argument required
        "https://www.uber.com/en/blog/berlin/engineering/rss",
        "https://github.blog/category/engineering/feed/",
        "https://medium.com/feed/miro-engineering"
    ]

    res = []
    for feed in feed_list:
        logging.info(f"Start process {feed}")
        parsed = feedparser.parse(feed)
        logging.info(f"Parsed {feed}")

        if len(parsed.entries) == 0:
            logging.error(f"Cannot parce {feed}")
            continue

        last_article = parsed.entries[0]
        title = last_article.title
        link = last_article.link

        article_published = datetime.fromtimestamp(mktime(last_article.published_parsed)).date()
        prev_day = datetime.today().date() - timedelta(days=1)

        if article_published == prev_day:
            res.append(f"{title} {link}")

    return "\n".join(res)


def blog_updates():
    try:
        return _blog_updates()
    except:
        logging.exception("Couldn't get blog posts")
        return "Couldn't get blog posts"
