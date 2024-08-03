package weather

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stringFormat(t *testing.T) {
	weatherData := WeatherData{"18.15  22.72  24.59", "Clear - Clouds", "21:18"}
	actual := weatherData.String()
	expected := `Weather today in Berlin Köpenick:
🌡️ 18.15  22.72  24.59
☀️ - ☁️
🌇 Sunset at 21:18`
	assert.Equal(t, expected, actual)
}