package repository

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const weatherResponse = `{
	"location": {
		"name": "Bauru",
		"region": "Sao Paulo",
		"country": "Brazil",
		"lat": -22.3167,
		"lon": -49.0667,
		"tz_id": "America/Sao_Paulo",
		"localtime_epoch": 1759111199,
		"localtime": "2025-09-28 22:59"
	},
	"current": {
		"last_updated_epoch": 1759110300,
		"last_updated": "2025-09-28 22:45",
		"temp_c": 25.1,
		"temp_f": 77.2,
		"is_day": 0,
		"condition": {
			"text": "Clear",
			"icon": "//cdn.weatherapi.com/weather/64x64/night/113.png",
			"code": 1000
		},
		"wind_mph": 4.0,
		"wind_kph": 6.5,
		"wind_degree": 158,
		"wind_dir": "SSE",
		"pressure_mb": 1017.0,
		"pressure_in": 30.03,
		"precip_mm": 0.0,
		"precip_in": 0.0,
		"humidity": 51,
		"cloud": 0,
		"feelslike_c": 25.2,
		"feelslike_f": 77.3,
		"windchill_c": 22.5,
		"windchill_f": 72.5,
		"heatindex_c": 23.9,
		"heatindex_f": 75.1,
		"dewpoint_c": 6.7,
		"dewpoint_f": 44.1,
		"vis_km": 10.0,
		"vis_miles": 6.0,
		"uv": 0.0,
		"gust_mph": 8.5,
		"gust_kph": 13.6,
		"short_rad": 0,
		"diff_rad": 0,
		"dni": 0,
		"gti": 0
	}
}`

func TestWeatherAPIRepository_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, weatherResponse)
	}))
	defer server.Close()

	repo := &WeatherAPIRepository{
		Key:     "mykey",
		BaseURL: server.URL,
	}

	result := repo.Get("Bauru")

	if result.Location.Name != "Bauru" {
		t.Errorf("expected location name 'Bauru', got '%s'", result.Location.Name)
	}

	if result.Current.TempC != 25.1 {
		t.Errorf("expected current temperature 25.1, got %f", result.Current.TempC)
	}

	if result.Current.Condition.Text != "Clear" {
		t.Errorf("expected current condition 'Clear', got '%s'", result.Current.Condition.Text)
	}
}
