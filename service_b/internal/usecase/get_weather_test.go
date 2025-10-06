package usecase_test

import (
	"testing"

	"github.com/dhquero/otel-go/service_b/internal/infra/repository"
	"github.com/dhquero/otel-go/service_b/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockViaCEPRepository struct {
	mock.Mock
}

func (m *MockViaCEPRepository) Get(zipCode string) repository.ViaCEP {
	args := m.Called(zipCode)
	return args.Get(0).(repository.ViaCEP)
}

type MockWeatherAPIRepository struct {
	mock.Mock
}

func (m *MockWeatherAPIRepository) Get(city string) repository.WeatherAPI {
	args := m.Called(city)
	return args.Get(0).(repository.WeatherAPI)
}

func TestGetWeatherUseCase_InvalidZipCode(t *testing.T) {
	mockViaCEP := new(MockViaCEPRepository)
	mockWeatherAPI := new(MockWeatherAPIRepository)

	useCase := usecase.NewGetWeatherUseCase(mockViaCEP, mockWeatherAPI)

	input := usecase.ZipCodeInputDTO{ZipCode: "123"}

	output, err := useCase.Get(input)

	assert.Nil(t, output)
	assert.EqualError(t, err, "invalid zipcode")
}

func TestGetWeatherUseCase_ZipCodeNotFound(t *testing.T) {
	mockViaCEP := new(MockViaCEPRepository)
	mockWeatherAPI := new(MockWeatherAPIRepository)

	useCase := usecase.NewGetWeatherUseCase(mockViaCEP, mockWeatherAPI)

	input := usecase.ZipCodeInputDTO{ZipCode: "98765432"}

	mockViaCEP.On("Get", "98765432").Return(repository.ViaCEP{Localidade: ""})

	output, err := useCase.Get(input)

	assert.Nil(t, output)
	assert.EqualError(t, err, "can not find zipcode")
	mockViaCEP.AssertExpectations(t)
}

func TestGetWeatherUseCase_Success(t *testing.T) {
	mockViaCEP := new(MockViaCEPRepository)
	mockWeatherAPI := new(MockWeatherAPIRepository)

	useCase := usecase.NewGetWeatherUseCase(mockViaCEP, mockWeatherAPI)

	input := usecase.ZipCodeInputDTO{ZipCode: "17010030"}

	mockViaCEP.On("Get", "17010030").Return(repository.ViaCEP{
		Cep:         "17010-030",
		Logradouro:  "Rua Presidente Kennedy",
		Complemento: "de Quadra 5 a Quadra 8",
		Unidade:     "",
		Bairro:      "Centro",
		Localidade:  "Bauru",
		Uf:          "SP",
		Estado:      "São Paulo",
		Regiao:      "Sudeste",
		Ibge:        "3506003",
		Gia:         "2094",
		Ddd:         "14",
		Siafi:       "6219",
	})

	mockWeatherAPI.On("Get", "Bauru").Return(repository.WeatherAPI{
		Location: repository.Location{
			Name:           "Bauru",
			Region:         "Sao Paulo",
			Country:        "Brazil",
			Lat:            -22.3167,
			Lon:            -49.0667,
			TzId:           "America/Sao_Paulo",
			LocaltimeEpoch: 1759111199,
			Localtime:      "2025-09-28 22:59",
		},
		Current: repository.Current{
			LastUpdatedEpoch: 1759110300,
			LastUpdated:      "2025-09-28 22:45",
			TempC:            25.1,
			TempF:            77.2,
			IsDay:            0,
			Condition: repository.Condition{
				Text: "Clear",
				Icon: "//cdn.weatherapi.com/weather/64x64/night/113.png",
				Code: 1000,
			},
			WindMph:    4.0,
			WindKph:    6.5,
			WindDegree: 158,
			WindDir:    "SSE",
			PressureMb: 1017.0,
			PressureIn: 30.03,
			PrecipMm:   0.0,
			PrecipIn:   0.0,
			Humidity:   51,
			Cloud:      0,
			FeelslikeC: 25.2,
			FeelslikeF: 77.3,
			WindchillC: 22.5,
			WindchillF: 72.5,
			HeatindexC: 23.9,
			HeatindexF: 75.1,
			DewpointC:  6.7,
			DewpointF:  44.1,
			VisKm:      10.0,
			VisMiles:   6.0,
			Uv:         0.0,
			GustMph:    8.5,
			GustKph:    13.6,
			ShortRad:   0,
			DiffRad:    0,
			DNI:        0,
			GTI:        0,
		},
	})

	output, err := useCase.Get(input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, 25.1, output.Celsius)
	assert.Equal(t, 77.18, output.Fahrenheit)
	assert.Equal(t, 298.1, output.Kelvin)

	mockViaCEP.AssertExpectations(t)
	mockWeatherAPI.AssertExpectations(t)
}
