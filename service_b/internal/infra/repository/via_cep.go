package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ViaCEPRepositoryInterface interface {
	Get(cep string) ViaCEP
}

type ViaCEPRepository struct {
	BaseURL string
}

func NewViaCEPRepository() *ViaCEPRepository {
	return &ViaCEPRepository{
		BaseURL: "http://viacep.com.br/ws",
	}
}

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func (v *ViaCEPRepository) Get(cep string) ViaCEP {
	req, err := http.Get(v.BaseURL + "/" + cep + "/json/")

	if err != nil {
		fmt.Fprintf(os.Stderr, "ViaCEP - Erro ao fazer requisição: %v\n", err)
	}

	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ViaCEP - Erro ao ler resposta: %v\n", err)
	}

	var data ViaCEP

	err = json.Unmarshal(res, &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ViaCEP - Erro ao fazer parse da resposta: %v\n", err)
	}

	return data
}
