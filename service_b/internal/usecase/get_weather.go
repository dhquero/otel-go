package usecase

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/dhquero/otel-go/service_b/internal/infra/repository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type ZipCodeInputDTO struct {
	ZipCode string
}

type TemperatureOutputDTO struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

type GetWeatherUseCase struct {
	ViaCEPRepositoryInterface     repository.ViaCEPRepositoryInterface
	WeatherAPIRepositoryInterface repository.WeatherAPIRepositoryInterface
	Tracer                        trace.Tracer
	request                       *http.Request
}

func NewGetWeatherUseCase(
	viaCEPRepository repository.ViaCEPRepositoryInterface,
	WeatherAPIRepositoryInterface repository.WeatherAPIRepositoryInterface,
	tracer trace.Tracer,
	request *http.Request,
) *GetWeatherUseCase {
	return &GetWeatherUseCase{
		ViaCEPRepositoryInterface:     viaCEPRepository,
		WeatherAPIRepositoryInterface: WeatherAPIRepositoryInterface,
		Tracer:                        tracer,
		request:                       request,
	}
}

func (z *ZipCodeInputDTO) validate() bool {
	match, _ := regexp.MatchString(`^\d{8}$`, z.ZipCode)
	return match
}

func (g *GetWeatherUseCase) Get(zipCodeInputDTO ZipCodeInputDTO) (*TemperatureOutputDTO, error) {
	if !zipCodeInputDTO.validate() {
		return nil, errors.New("invalid zipcode")
	}

	carrier := propagation.HeaderCarrier(g.request.Header)
	ctx := g.request.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	_, span_via_cep_api := g.Tracer.Start(ctx, "via_cep_api")
	viaCEP := g.ViaCEPRepositoryInterface.Get(zipCodeInputDTO.ZipCode)
	if viaCEP.Localidade == "" {
		return nil, errors.New("can not find zipcode")
	}
	span_via_cep_api.End()

	_, span_weather_api := g.Tracer.Start(ctx, "weather_api")
	weatherAPI := g.WeatherAPIRepositoryInterface.Get(viaCEP.Localidade)
	span_weather_api.End()

	return &TemperatureOutputDTO{
		City:       viaCEP.Localidade,
		Celsius:    weatherAPI.Current.TempC,
		Fahrenheit: weatherAPI.Current.TempC*1.8 + 32,
		Kelvin:     weatherAPI.Current.TempC + 273,
	}, nil
}
