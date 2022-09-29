import logging

from decouple import config
import telepot

from modules.blogs import blog_updates
from modules.currency import curencies
from modules.herthaTickets import hertha_tickets
from modules.weather import weather

TELEGRAM_TO = config("TELEGRAM_TO")
TELEGRAM_TOKEN = config("TELEGRAM_TOKEN")
OPEN_WHEATHER_API_KEY = config("OPEN_WHEATHER_API_KEY")
EXCHANGERATE_API_KEY = config("EXCHANGERATE_API_KEY")

logging.basicConfig(level=logging.INFO,
                    format="%(asctime)s,%(msecs)d %(name)s %(levelname)s %(message)s",
                    filemode="w")

def send_telegram_message(msg):
    bot = telepot.Bot(TELEGRAM_TOKEN)
    bot.getMe()
    bot.sendMessage(TELEGRAM_TO, msg)


if __name__ == "__main__":
    metrics = [
        weather(OPEN_WHEATHER_API_KEY),
        curencies(EXCHANGERATE_API_KEY),
        blog_updates(),
        hertha_tickets()
    ]
    res_msg = "\n\n".join(metrics)
    send_telegram_message(res_msg)
