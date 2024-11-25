package weather

import (
	"fmt"
	"strings"
)

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
	precipitationEmoji := replacePrecipitationToEmoji(w.Precipitation)
	weatherMsg += fmt.Sprintf("🌡️ %s\n%s\n🌇 Sunset at %s", w.Temp, precipitationEmoji, w.SunsetTime)
	return weatherMsg
}

func replacePrecipitationToEmoji(s string) string {
	s = strings.ReplaceAll(s, "Clouds", "☁️")
	s = strings.ReplaceAll(s, "Clear", "☀️")
	s = strings.ReplaceAll(s, "Rain", "🌧️")
	s = strings.ReplaceAll(s, "Snow", "❄️")
	return s
}
