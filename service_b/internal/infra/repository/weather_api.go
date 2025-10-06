package repository

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type WeatherAPIRepositoryInterface interface {
	Get(city string) WeatherAPI
}

type WeatherAPIRepository struct {
	Key     string
	BaseURL string
}

func NewWeatherAPIRepository(key string) *WeatherAPIRepository {
	return &WeatherAPIRepository{
		Key:     key,
		BaseURL: "https://api.weatherapi.com/v1/current.json",
	}
}

type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzId           string  `json:"tz_id"`
	LocaltimeEpoch int     `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

type Current struct {
	LastUpdatedEpoch int       `json:"last_updated_epoch"`
	LastUpdated      string    `json:"last_updated"`
	TempC            float64   `json:"temp_c"`
	TempF            float64   `json:"temp_f"`
	IsDay            int       `json:"is_day"`
	Condition        Condition `json:"condition"`
	Localtime        string    `json:"localtime"`
	WindMph          float64   `json:"wind_mph"`
	WindKph          float64   `json:"wind_kph"`
	WindDegree       int       `json:"wind_degree"`
	WindDir          string    `json:"wind_dir"`
	PressureMb       float64   `json:"pressure_mb"`
	PressureIn       float64   `json:"pressure_in"`
	PrecipMm         float64   `json:"precip_mm"`
	PrecipIn         float64   `json:"precip_in"`
	Humidity         int       `json:"humidity"`
	Cloud            int       `json:"cloud"`
	FeelslikeC       float64   `json:"feelslike_c"`
	FeelslikeF       float64   `json:"feelslike_f"`
	WindchillC       float64   `json:"windchill_c"`
	WindchillF       float64   `json:"windchill_f"`
	HeatindexC       float64   `json:"heatindex_c"`
	HeatindexF       float64   `json:"heatindex_f"`
	DewpointC        float64   `json:"dewpoint_c"`
	DewpointF        float64   `json:"dewpoint_f"`
	VisKm            float64   `json:"vis_km"`
	VisMiles         float64   `json:"vis_miles"`
	Uv               float64   `json:"uv"`
	GustMph          float64   `json:"gust_mph"`
	GustKph          float64   `json:"gust_kph"`
	ShortRad         float64   `json:"short_rad"`
	DiffRad          float64   `json:"diff_rad"`
	DNI              float64   `json:"dni"`
	GTI              float64   `json:"gti"`
}

type WeatherAPI struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

func (w *WeatherAPIRepository) Get(city string) WeatherAPI {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := httpClient.Get(w.BaseURL + "?q=" + city + "&key=" + w.Key)

	if err != nil {
		fmt.Fprintf(os.Stderr, "WeatherAPI - Erro ao fazer requisição: %v\n", err)
	}

	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "WeatherAPI - Erro ao ler resposta: %v\n", err)
	}

	var data WeatherAPI

	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WeatherAPI - Erro ao fazer parse da resposta: %v\n", err)
	}

	return data
}
