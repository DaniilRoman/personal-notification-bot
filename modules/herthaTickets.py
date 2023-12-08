import requests
from bs4 import BeautifulSoup
import logging

from utils.dynamodb import ItemStoreService


def __hertha_tickets():
    res = ""
    hertha_tickets_url = "https://ticket-onlineshop.com/ols/hbsctk/en/tk/"
    headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0'}
    html_text = requests.get(hertha_tickets_url, timeout=5, headers=headers).text
    soup = BeautifulSoup(html_text, 'html.parser')

    for event in soup.findAll(True, {"class": ['event-card__headings']}):
        res += event.text.replace("\n", "").replace("   ", "").replace("Hertha BSC", "Hertha BSC -")
        res += "\n"

    return res


def hertha_tickets(item_store_service: ItemStoreService) -> str:
    try:
        actual_tickets = __hertha_tickets()
        if item_store_service is None:
            return actual_tickets
        stored_tickets = item_store_service.get_item("hertha_tickets")
        if actual_tickets != stored_tickets:
            item_store_service.save_item("hertha_tickets", actual_tickets)
            return actual_tickets
        else:
            return ""
    except:
        logging.exception("Couldn't get Hertha tickets")
        return "Couldn't get Hertha tickets"
