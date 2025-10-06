package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dhquero/otel-go/pkg/otel"
	"github.com/dhquero/otel-go/service_b/internal/infra/web"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	gootel "go.opentelemetry.io/otel"
)

func main() {
	otel.InitTracer(os.Getenv("ZIPKIN_URL"), "service_b")

	mux := http.NewServeMux()

	tracer := gootel.Tracer("service_b")
	handler := web.NewHandlerWeather(tracer)

	mux.Handle("/", otelhttp.NewHandler(http.HandlerFunc(handler.HandleWeather), "weather"))

	fmt.Println("Server starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
