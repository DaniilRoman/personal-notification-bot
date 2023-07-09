import logging

from typing import List

from modules.blogs import blog_updates
from modules.currency import curencies, CurrencyData
from modules.herthaTickets import hertha_tickets
from modules.weather import weather
from utils.chatgpt_summarizing import configure_openai
from utils.dynamodb import DynamodbConfig, ItemStoreService
from utils.templating import render_index_html

logging.basicConfig(level=logging.INFO,
                    format="%(asctime)s,%(msecs)d %(name)s %(levelname)s %(message)s",
                    filemode="w")

dynamodb_config = DynamodbConfig(endpoint_url="http://localhost:8000")
item_store_service = ItemStoreService(dynamodb_config)
configure_openai("", "")


def _send_telegram_message(data_list):
    msg = _create_telegram_message(data_list)
    print(msg)


def _create_telegram_message(data_list: List[object]) -> str:
    str_list = [str(data) for data in data_list]
    return "\n\n".join(str_list)


if __name__ == "__main__":
    currency_data = CurrencyData()
    currency_data.key_values["EUR"] = 97.01
    weather_data = "Weather good"
    blogs_data = blog_updates()
    hertha_tickets_data = hertha_tickets(item_store_service)

    _send_telegram_message([weather_data, currency_data, blogs_data, hertha_tickets_data])

    data = {
        "weather": weather_data,
        "currency": currency_data,
        "blogs": blogs_data,
        "herta_tickets": hertha_tickets_data
    }
    render_index_html(data)
