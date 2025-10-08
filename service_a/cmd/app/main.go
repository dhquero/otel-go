package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dhquero/otel-go/pkg/otel"
	"github.com/dhquero/otel-go/service_a/internal/infra/web"
)

func main() {
	otel.InitTracer(os.Getenv("OTEL_COLLECTOR_URL"), "service_a")

	mux := http.NewServeMux()

	mux.HandleFunc("/", web.HandleWeather)

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
