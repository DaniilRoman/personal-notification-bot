import logging
from typing import List

from decouple import config
import telepot

from modules.blogs import blog_updates
from modules.currency import curencies
from modules.herthaTickets import hertha_tickets
from modules.unionBerlinTickets import union_berlin_tickets
from modules.weather import weather
from modules.mobileNumber import mobile_number_notification
from modules.wordOfTheDay import word_of_the_day
from utils.chatgpt_summarizing import configure_openai
from utils.dynamodb import DynamodbConfig, ItemStoreService
from utils.templating import render_index_html

TELEGRAM_TO = config("TELEGRAM_TO")
TELEGRAM_TOKEN = config("TELEGRAM_TOKEN")

OPEN_WHEATHER_API_KEY = config("OPEN_WHEATHER_API_KEY")
EXCHANGERATE_API_KEY = config("EXCHANGERATE_API_KEY")

AWS_ACCESS_KEY_ID = config("AWS_ACCESS_KEY_ID")
AWS_SECRET_ACCESS_KEY = config("AWS_SECRET_ACCESS_KEY")
REGION_NAME = config("REGION_NAME")

OPENAI_ACCESS_KEY = config("OPENAI_ACCESS_KEY")
OPENAI_ORGANIZATION = config("OPENAI_ORGANIZATION")

logging.basicConfig(level=logging.INFO,
                    format="%(asctime)s,%(msecs)d %(name)s %(levelname)s %(message)s",
                    filemode="w")

dynamodb_config = DynamodbConfig(AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, REGION_NAME)
item_store_service = ItemStoreService(dynamodb_config)
configure_openai(OPENAI_ACCESS_KEY, OPENAI_ORGANIZATION)


def _send_telegram_message(data_list):
    msg = _create_telegram_message(data_list)
    bot = telepot.Bot(TELEGRAM_TOKEN)
    bot.getMe()
    bot.sendMessage(TELEGRAM_TO, msg, parse_mode='Markdown')


def _create_telegram_message(data_list: List[object]) -> str:
    str_list = [str(data) for data in data_list]
    str_list.append("[Html page](https://daniilroman.github.io/personal-notification-bot/)")
    return "\n\n".join(str_list)


if __name__ == "__main__":
    weather_data = weather(OPEN_WHEATHER_API_KEY)
    currency_data = curencies(EXCHANGERATE_API_KEY)
    blogs_data = blog_updates()
    hertha_tickets_data = hertha_tickets(item_store_service)
    union_berlin_tickets_data = union_berlin_tickets(item_store_service)
    mobile_number_notification_data = mobile_number_notification()
    word_of_the_day_data = word_of_the_day()

    _send_telegram_message([weather_data, currency_data, blogs_data, hertha_tickets_data, union_berlin_tickets_data, mobile_number_notification_data, word_of_the_day_data])

    data = {
        "weather": weather_data,
        "currency": currency_data,
        "blogs": blogs_data,
        "herta_tickets": hertha_tickets_data
    }
    render_index_html(data)
