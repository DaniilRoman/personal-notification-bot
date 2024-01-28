package weather

import (
	"fmt"
	"log"
	"main/utils"
	"strings"
	"time"
)

const weatherUrl = "https://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&cnt=6&units=metric&appid=%s"


func GetWeather(token string) *WeatherData {
	res, err := getWeather(token)
    if err != nil {
		log.Printf("Error in Weather module: %s", err)
	 }
	return res
}

func getWeather(token string) (*WeatherData, error) {
	// Berlin Koepenick
	lat := 52.4514534
	lon := 13.5699097
	url := fmt.Sprintf(weatherUrl, lat, lon, token)

	weatherResponse := weatherResponse{}
	err := utils.DoGet(url, &weatherResponse)
	if err != nil {
		return nil, err
	}

	return weatherResponse.GetTodayWeather(), nil
}


func (w weatherResponse) GetTodayWeather() *WeatherData {
	tempNow := w.List[1].Main.Temp // +3 hours ~09:00
	tempAfternoon := w.List[2].Main.Temp // +6 hours ~12:00
	tempEvening := w.List[4].Main.Temp // +12 hours ~18:00

	tempStr := fmt.Sprintf("%.2f  %.2f  %.2f", tempNow, tempAfternoon, tempEvening)

	var precipitation []string
	for _, entry := range w.List {
		for _, weather := range entry.Weather {
			precipitation = append(precipitation, weather.Main)
		}
	}
	precipitation = utils.RemoveDuplicate(precipitation)
	precipitationStr := strings.Join(precipitation, " - ")

	sunsetTime := getSunset(w)

	return &WeatherData{tempStr, precipitationStr, sunsetTime}
}

func getSunset(w weatherResponse) string {
	timeUtc := time.Unix(w.City.Sunset, 0)

	timezone := "Europe/Berlin"
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return ""
	}

	timeWithTimezone := timeUtc.In(loc)
	return timeWithTimezone.Format("15:04")
}


type weatherResponse struct {
	City struct {
		Name  string `json:"name"`
		Sunset int64  `json:"sunset"`
	} `json:"city"`
	List []struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Main string `json:"main"`
		} `json:"weather"`
	} `json:"list"`
}
