import requests
import logging


def _curencies(EXCHANGERATE_API_KEY):
    currencies_msg = "Currencies:\n"
    currencies = ["USD", "EUR"]
    for c in currencies:
        url = f'https://v6.exchangerate-api.com/v6/{EXCHANGERATE_API_KEY}/latest/{c}'

        response = requests.get(url)
        data = response.json()
        currencies_msg += f'{c}: {data["conversion_rates"]["RUB"]} RUB\n'
    return currencies_msg


def curencies(token):
    try:
        return _curencies(token)
    except:
        logging.exception("Couldn't get currency")
        return "Couldn't get currency"
