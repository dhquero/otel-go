package repository

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const CEPResponse = `{
  "cep": "17010-030",
  "logradouro": "Rua Presidente Kennedy",
  "complemento": "de Quadra 5 a Quadra 8",
  "unidade": "",
  "bairro": "Centro",
  "localidade": "Bauru",
  "uf": "SP",
  "estado": "São Paulo",
  "regiao": "Sudeste",
  "ibge": "3506003",
  "gia": "2094",
  "ddd": "14",
  "siafi": "6219"
}`

func TestViaCEPRepository_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, CEPResponse)
	}))
	defer server.Close()

	repo := &ViaCEPRepository{
		BaseURL: server.URL,
	}

	result := repo.Get("17010030")

	if result.Cep != "17010-030" {
		t.Errorf("cep expected '17010-030', got '%s'", result.Cep)
	}

	if result.Bairro != "Centro" {
		t.Errorf("expected 'Centro', got '%s'", result.Bairro)
	}

	if result.Localidade != "Bauru" {
		t.Errorf("expected 'Bauru', got '%s'", result.Localidade)
	}
}
