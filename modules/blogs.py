from datetime import datetime, timedelta
from time import mktime

import logging
import feedparser


def _getBlogUpdates():
    feed_list = [
        "https://vladmihalcea.com/blog/feed/",
        "https://piotrminkowski.com/feed/",
        "https://blog.alexellis.io/rss/",
        "https://slack.engineering/feed/",
        "https://engineering.fb.com/feed/",
        "https://engineering.atspotify.com/feed/",
        "https://engineering.zalando.com/atom.xml"
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

def getBlogUpdates():
    try:
        return _getBlogUpdates()
    except:
        logging.exception("Couldn't get blog posts")
        return "Couldn't get blog posts"