package usecase

import (
	"errors"
	"regexp"

	"github.com/dhquero/otel-go/service_a/internal/infra/repository"
)

type CepInputDTO struct {
	Cep string `json:"cep"`
}

type TemperatureOutputDTO struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

type GetWeatherUseCase struct {
	ServiceBRepository repository.ServiceBRepositoryInterface
}

func NewGetWeatherUseCase(
	ServiceBRepository repository.ServiceBRepositoryInterface,
) *GetWeatherUseCase {
	return &GetWeatherUseCase{
		ServiceBRepository: ServiceBRepository,
	}
}

func (z *CepInputDTO) validate() bool {
	match, _ := regexp.MatchString(`^\d{8}$`, z.Cep)
	return match
}

func (g *GetWeatherUseCase) Get(zipCodeInputDTO CepInputDTO) (*TemperatureOutputDTO, error) {
	if !zipCodeInputDTO.validate() {
		return nil, errors.New("invalid zipcode")
	}

	weatherAPI := g.ServiceBRepository.Get(zipCodeInputDTO.Cep)

	if weatherAPI.City == "" {
		return nil, errors.New("can not find zipcode")
	}

	return &TemperatureOutputDTO{
		City:       weatherAPI.City,
		Celsius:    weatherAPI.Celsius,
		Fahrenheit: weatherAPI.Fahrenheit,
		Kelvin:     weatherAPI.Kelvin,
	}, nil
}
