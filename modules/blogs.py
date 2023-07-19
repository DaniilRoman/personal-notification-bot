from datetime import datetime, timedelta
from time import mktime

import logging
from typing import List

from bs4 import BeautifulSoup
import feedparser
import requests

from utils.chatgpt_summarizing import summarize_text

__blacklist_labels = ["android", "ios", "redux", "react", "frontend", "ui/ux", "career stories", "meeting", "spotlight",
                      "internship", "javascript", "css", "html", "typescript", "mobile", "uikit", "интерфейс", "дизайн",
                      "мобильн", "design", "interface", "A Bootiful Podcast", "Spark", "Docker Desktop"]

__feed_list = [
    # Personal blogs
    "https://vladmihalcea.com/blog/feed/",
    "https://piotrminkowski.com/feed/",
    "https://blog.alexellis.io/rss/",

    # Interesting company blogs
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
    "https://stackoverflow.blog/engineering/feed/",
    "https://medium.com/feed/tinder",
    "https://medium.com/feed/bbc-product-technology",
    "https://open.nytimes.com/feed",
    "https://neo4j.com/developer-blog/feed/",

    # Research labs
    # https://githubnext.com, # didn't find rss feed,
    "http://feeds.feedburner.com/blogspot/gJZg",  # Google research
    "https://openai.com/blog/rss.xml",
    "https://research.facebook.com/feed/",
    "https://www.microsoft.com/en-us/research/feed/",
    "https://bair.berkeley.edu/blog/feed",

    # Tech
    "https://circleci.com/blog/feed.xml",
    "https://kubernetes.io/feed.xml",
    "https://www.docker.com/blog/feed",
    "https://redis.com/blog/rss",
    "https://www.mongodb.com/blog/rss",
    "https://debezium.io/blog.atom",
    "https://www.elastic.co/blog/feed",
    "https://aws.amazon.com/blogs/database/tag/dynamodb/feed",
    "https://about.gitlab.com/atom.xml",
    "https://feeds.feedburner.com/ContinuousBlog/",  # Jenkins
    "https://in.relation.to/blog.atom",  # Hibernate
    "https://www.cncf.io/blog/feed/",

    # Others company bogs
    "https://www.hashicorp.com/blog/feed.xml",
    "https://microservices.io/feed.xml",
    "https://www.confluent.io/rss.xml",
    "https://blog.cloudflare.com/rss",
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
    "https://www.reddit.com/r/RedditEng/.rss",

    # Company blogs to delete
    "https://engineering.fb.com/feed/",
    "https://blog.twitter.com/engineering/en_us/blog.rss",
    "https://www.datadoghq.com/blog/index.xml",
    "https://grafana.com/categories/engineering/index.xml"
]


class BlogData:
    def __init__(self, link, title, img="", summary=""):
        self.link = link
        self.img = img
        self.title = title
        self.summary = summary

    def __repr__(self):
        website_name = self.link.replace("https://", "").replace("http://", "").split("/")[0]
        res_article_str = f"- [{self.title}]({self.link})\n[[{website_name}]]"
        return res_article_str


class BlogsData:
    def __init__(self):
        self.blogs: List[BlogData] = []

    def append(self, blog_data: BlogData):
        self.blogs.append(blog_data)

    def __repr__(self):
        return "\n".join([str(blog) for blog in self.blogs])


def __none_blacklist_labels(parsed_article):
    if any(black_list_label in parsed_article.title.lower() for black_list_label in __blacklist_labels):
        return False
    if parsed_article.get("tags") is not None:
        for tag in parsed_article.tags:
            if any(black_list_label == tag.term.lower() for black_list_label in __blacklist_labels):
                return False

    return True


def _set_extra_fields(blog_data):
    try:
        page = requests.get(blog_data.link)
        soup = BeautifulSoup(page.text, "html.parser")

        img = _get_img(soup, blog_data.link)
        blog_data.img = img

        text_to_summarize = soup.getText().replace("\n", " ").replace("\t", " ").replace("  ", " ")
        blog_data.summary = summarize_text(text_to_summarize)
    except Exception as ex:
        logging.warning(f"Couldn't get img or summarize text for {blog_data.link}", ex)
        return ""


def __add_new_article_to_res_list(feed, blogs_data: BlogsData):
    logging.info(f"Start process {feed}")
    parsed = feedparser.parse(feed)
    logging.info(f"Parsed {feed}")

    if len(parsed.entries) == 0:
        logging.error(f"Cannot parce {feed}")
        return

    last_article = parsed.entries[0]

    if _is_article_published_yesterday(last_article, parsed):

        blog_data = BlogData(last_article.link, last_article.title)
        if __none_blacklist_labels(last_article):
            _set_extra_fields(blog_data)
            blogs_data.append(blog_data)
        else:
            logging.info(f"Filtered by topic: {str(blog_data)}")


def _is_article_published_yesterday(last_article, parsed) -> bool:
    if last_article.get("published_parsed") is None:
        if parsed.get("updated_parsed") is None:
            article_published = datetime.fromtimestamp(mktime(last_article.updated)).date()
        else:
            article_published = datetime.fromtimestamp(mktime(parsed.updated_parsed)).date()
    else:
        article_published = datetime.fromtimestamp(mktime(last_article.published_parsed)).date()
    prev_day = datetime.today().date() - timedelta(days=1)
    return article_published == prev_day


def _get_img(soup, link):
    try:
        img = soup.find("meta", property="og:image").attrs["content"]
        if img is None:
            return ""
        return img
    except Exception as ex:
        logging.warning(f"Couldn't get img for {link}", ex)
        return ""


def _blog_updates() -> BlogsData:
    res = BlogsData()
    for feed in __feed_list:
        try:
            __add_new_article_to_res_list(feed, res)
        except:
            logging.exception(f"Couldn't process blog: {str(feed)}")

    return res


def blog_updates() -> BlogsData:
    try:
        return _blog_updates()
    except:
        logging.exception("Couldn't get blog posts")
        return "Couldn't get blog posts"  # TODO fix
