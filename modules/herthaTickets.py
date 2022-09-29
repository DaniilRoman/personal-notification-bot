import requests
from bs4 import BeautifulSoup
import logging


def _hertha_tickets():
    res = ""
    hertha_tickets_url = "https://ticket-onlineshop.com/ols/hbsctk/en/tk/"
    headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0'}
    html_text = requests.get(hertha_tickets_url, timeout=5, headers=headers).text
    soup = BeautifulSoup(html_text, 'html.parser')

    for event in soup.findAll(True, {"class": ['event-card__headings']}):
        res += event.text.replace("\n", "").replace("   ", "").replace("Hertha BSC", "Hertha BSC -")
        res += "\n"

    return res


def hertha_tickets():
    try:
        return _hertha_tickets()
    except:
        logging.exception("Couldn't get Hertha tickets")
        return "Couldn't get Hertha tickets"
