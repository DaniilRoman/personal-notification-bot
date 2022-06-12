from decouple import config
import requests
import telepot
from yaweather import YaWeather, Russia

TELEGRAM_TO = config("TELEGRAM_TO")
TELEGRAM_TOKEN = config("TELEGRAM_TOKEN")
YANDEX_WHETHER_KEY = config("YANDEX_WHETHER_KEY")
EXCHANGERATE_API_KEY = config("EXCHANGERATE_API_KEY")


def getWeather():
    y = YaWeather(api_key=YANDEX_WHETHER_KEY)
    res = y.forecast(Russia.NizhniyNovgorod)

    weather = f'Now: {res.fact.temp} °C, feels like {res.fact.feels_like} °C'

    rain = ""
    conditions = [f.condition for f in res.forecasts]
    if 'rain' in conditions:
        rain = "Will be rain."
    return f'Weather:\n{weather}\n{rain}'

def getCurencies():
    currencies_msg = "Currencies:\n"
    currencies = ["USD", "EUR"]
    for c in currencies:
        url = f'https://v6.exchangerate-api.com/v6/{EXCHANGERATE_API_KEY}/latest/{c}'

        response = requests.get(url)
        data = response.json()
        currencies_msg += f'{c}: {data["conversion_rates"]["RUB"]} RUB\n'
    return currencies_msg


def send_telegram_message(msg):
    bot = telepot.Bot(TELEGRAM_TOKEN)
    bot.getMe()
    bot.sendMessage(TELEGRAM_TO, msg)


if __name__ == "__main__":
    msg = getWeather()
    msg += "\n"
    msg += getCurencies()
    print(msg)
    send_telegram_message(msg)
