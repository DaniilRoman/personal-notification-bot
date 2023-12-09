package weather

import "fmt"

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
