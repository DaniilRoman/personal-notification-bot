import requests
from bs4 import BeautifulSoup
import logging

from utils.dynamodb import ItemStoreService


def __union_berlin_tickets():
    union_berlin_tickets_url = "https://tickets.union-zeughaus.de/unveu/heimspiele_2.htm"
    res = f"[Union Berlin tickets]({union_berlin_tickets_url}):\n"
    headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:82.0) Gecko/20100101 Firefox/82.0'}
    html_text = requests.get(union_berlin_tickets_url, timeout=5, headers=headers).text
    soup = BeautifulSoup(html_text, 'html.parser')

    for event in soup.findAll(True, {"class": ['ticket listitem gamehome']}):
        res += " - ".join([team.text for team in event.find_all_next("h2")])
        res += "\n"

    return res


def union_berlin_tickets(item_store_service: ItemStoreService) -> str:
    try:
        actual_tickets = __union_berlin_tickets()
        if item_store_service is None:
            return actual_tickets
        stored_tickets = item_store_service.get_item("union_berlin_tickets")
        if actual_tickets != stored_tickets:
            item_store_service.save_item("union_berlin_tickets", actual_tickets)
            return actual_tickets
        else:
            return ""
    except:
        logging.exception("Couldn't get Union Berlin tickets")
        return "Couldn't get Union Berlin tickets"
