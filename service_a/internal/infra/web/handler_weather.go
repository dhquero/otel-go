package web

import (
	"encoding/json"
	"net/http"

	"github.com/dhquero/otel-go/service_a/internal/infra/repository"
	"github.com/dhquero/otel-go/service_a/internal/usecase"
)

func HandleWeather(w http.ResponseWriter, r *http.Request) {
	var zipCodeInputDTO usecase.CepInputDTO

	err := json.NewDecoder(r.Body).Decode(&zipCodeInputDTO)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	serviceBRepository := repository.NewServiceBRepository()

	getWeatherUseCase := usecase.NewGetWeatherUseCase(serviceBRepository)

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
