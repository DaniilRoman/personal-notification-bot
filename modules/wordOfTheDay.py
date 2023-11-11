import requests
from bs4 import BeautifulSoup
import logging


def __word_of_the_day():
    union_berlin_tickets_url = "https://www.nytimes.com/column/learning-word-of-the-day"
    headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0'}
    html_text = requests.get(union_berlin_tickets_url, timeout=5, headers=headers).text
    soup = BeautifulSoup(html_text, 'html.parser')

    last_article = soup.find('article')
    last_article_url = "https://www.nytimes.com" + last_article.find('a')['href']
    last_article_title = last_article.find('h3').text

    res =f"[{last_article_title}]({last_article_url}):\n"

    return res


def word_of_the_day() -> str:
    try:
        return __word_of_the_day()
    except:
        logging.exception("Couldn't get word of the day")
        return "Couldn't get word of the day"
