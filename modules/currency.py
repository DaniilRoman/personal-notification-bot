import requests
import logging


class CurrencyData:
    def __init__(self):
        self.key_values = {}

    def __repr__(self):
        if self.key_values == {}:
            return "Couldn't get currency"
        currencies_msg = "Currencies:\n"
        for k, v in self.key_values.items():
            currencies_msg += f'{k}: {v} RUB\n'
        return currencies_msg


def _curencies(EXCHANGERATE_API_KEY):
    currency_data = CurrencyData()
    currencies = ["USD", "EUR"]
    for c in currencies:
        url = f'https://v6.exchangerate-api.com/v6/{EXCHANGERATE_API_KEY}/latest/{c}'

        response = requests.get(url)
        data = response.json()
        currency_data.key_values[c] = data["conversion_rates"]["RUB"]
    return currency_data


def curencies(token) -> CurrencyData:
    try:
        return _curencies(token)
    except:
        logging.exception("Couldn't get currency")
        return CurrencyData()
