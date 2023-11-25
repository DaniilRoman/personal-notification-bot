package modules

import (
	"fmt"
	"log"
	utils "main/utils"
	"strings"
	"time"
)


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
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&cnt=6&units=metric&appid=%s", lat, lon, token)

	weatherResponse := weatherResponse{}
	err := utils.DoGet(url, &weatherResponse)
	if err != nil {
		return nil, err
	}

	return weatherResponse.GetTodayWeather(), nil
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

type WeatherData struct {
	Temp string
	Precipitation string
	SunsetTime string
}

func (w *WeatherData) String() string {
	if w == nil {
		return ""
	}
	weatherMsg := "Weather today in Berlin Köpenick:\n"
	weatherMsg += fmt.Sprintf("%s\n%s\nSunset at %s", w.Temp, w.Precipitation, w.SunsetTime)
	return weatherMsg
}

func (w weatherResponse) GetTodayWeather() *WeatherData {
	tempNow := w.List[1].Main.Temp // +3 hours ~10:00
	tempAfternoon := w.List[2].Main.Temp // +6 hours ~13:00
	tempEvening := w.List[4].Main.Temp // +12 hours ~19:00

	tempStr := fmt.Sprintf("%.2f : %.2f : %.2f", tempNow, tempAfternoon, tempEvening)

	var precipitation []string
	for _, entry := range w.List {
		for _, weather := range entry.Weather {
			precipitation = append(precipitation, weather.Main)
		}
	}
	precipitation = utils.RemoveDuplicate(precipitation)
	precipitationStr := strings.Join(precipitation, " - ")

	sunsetTime := time.Unix(w.City.Sunset, 0).Format("15:04")

	return &WeatherData{tempStr, precipitationStr, sunsetTime}
}