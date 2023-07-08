from datetime import datetime
import requests
import logging


def _timestamp_to_date(timestamp):
    return datetime.fromtimestamp(timestamp).time()


def _weather(token) -> str:
    # Berlin Koepenick
    lat = 52.4514534
    lon = 13.5699097
    url = f'https://api.openweathermap.org/data/2.5/forecast?lat={lat}&lon={lon}&cnt=6&units=metric&appid={token}'

    response = requests.get(url)
    weather = response.json()

    city_name = weather["city"]["name"]
    sunset_time = _timestamp_to_date(weather["city"]["sunset"])

    temp_now = weather["list"][1]["main"]["temp"] # +3 hours ~10:00
    temp_afternoon = weather["list"][2]["main"]["temp"]  # +6 hours ~13:00
    temp_evening = weather["list"][4]["main"]["temp"]  # +12 hours ~19:00

    temp_str = f"{temp_now} :: {temp_afternoon} :: {temp_evening}"

    precipitation = " - ".join(set(w["main"] for entry in weather["list"] for w in entry["weather"]))

    return f"Today in {city_name}: {temp_str} \nWeather: {precipitation} \nSunset at {sunset_time}"


def weather(token) -> str:
    try:
        return _weather(token)
    except:
        logging.exception("Couldn't get weather")
        return "Couldn't get weather"
