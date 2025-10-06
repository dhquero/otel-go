package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type ServiceBRepositoryInterface interface {
	Get(city string) ServiceB
}

type ServiceBRepository struct {
	BaseURL string
}

func NewServiceBRepository() *ServiceBRepository {
	return &ServiceBRepository{
		BaseURL: "http://service-b:8081/",
	}
}

type ServiceB struct {
	City       string  `json:"city"`
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func (w *ServiceBRepository) Get(cep string) ServiceB {
	ctx := context.Background()

	request, err := http.NewRequestWithContext(ctx, "GET", w.BaseURL+"?cep="+cep, nil)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ServiceB - Fala ao criar requisição: %v\n", err)
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))

	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithSpanNameFormatter(func(_ string, request *http.Request) string {
				return "get-service-b"
			}),
		),
	}

	req, err := httpClient.Do(request)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ServiceB - Erro ao fazer requisição: %v\n", err)
	}

	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ServiceB - Erro ao ler resposta: %v\n", err)
	}

	var data ServiceB

	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ServiceB - Erro ao fazer parse da resposta: %v\n", err)
	}

	return data
}
