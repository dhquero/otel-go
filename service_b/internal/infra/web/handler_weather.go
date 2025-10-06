package web

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/dhquero/otel-go/service_b/internal/infra/repository"
	"github.com/dhquero/otel-go/service_b/internal/usecase"
	"go.opentelemetry.io/otel/trace"
)

type HandlerWeather struct {
	Tracer trace.Tracer
}

func NewHandlerWeather(tracer trace.Tracer) *HandlerWeather {
	return &HandlerWeather{
		Tracer: tracer,
	}
}

func (h *HandlerWeather) HandleWeather(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	cep := queryParams.Get("cep")

	viaCEPrepository := repository.NewViaCEPRepository()
	weatherAPIRepository := repository.NewWeatherAPIRepository(os.Getenv("WEATHER_API_KEY"))

	getWeatherUseCase := usecase.NewGetWeatherUseCase(viaCEPrepository, weatherAPIRepository, h.Tracer, r)

	zipCodeInputDTO := usecase.ZipCodeInputDTO{
		ZipCode: cep,
	}

	temperature, err := getWeatherUseCase.Get(zipCodeInputDTO)

	if err != nil {
		errorMessage := err.Error()

		if errorMessage == "invalid zipcode" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(err.Error()))
			return
		}

		if errorMessage == "can not find zipcode" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(temperature)
}
