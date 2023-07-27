package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

type StringerInterface interface {
	ToString() string
}
type ViaCepResultDTO struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type ApicepResultDTO struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func (p ViaCepResultDTO) ToString() string {
	return stringfy(p)
}

func (c ApicepResultDTO) ToString() string {
	return stringfy(c)
}

func stringfy(dto interface{}) string {
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return fmt.Sprintf("Erro ao converter para JSON: %v", err)
	}
	return string(jsonBytes)
}

type MessageReturn struct {
	Message  string `json:"response-time"`
	Response string `json:"response"`
}

type ApiReturn struct {
	Message string `json:"message"`
	Body    []byte `json:"response"`
}

func main() {

	cep := flag.String("cep", "", "URL da primeira API")
	flag.Parse()

	urlAPI1 := "https://cdn.apicep.com/file/apicep/" + *cep + ".json"
	urlAPI2 := "http://viacep.com.br/ws/" + *cep + "/json/"
	timeout := 1 * time.Second

	ch := make(chan interface{})
	ch2 := make(chan interface{})

	apiCep := ApicepResultDTO{}
	go workerCdnApicep(ch, urlAPI1, apiCep)

	viaCepDto := ViaCepResultDTO{}
	go workerViaCep(ch2, urlAPI2, viaCepDto)

	select {
	case result := <-ch:
		out, _ := json.Marshal(result)
		fmt.Println(string(out))
	case result := <-ch2:
		out, _ := json.Marshal(result)
		fmt.Println(string(out))
	case <-time.After(timeout):
		fmt.Println("Erro: timeout ao aguardar a resposta.")
		return
	}

}

func workerViaCep(ch chan<- interface{}, url string, dto ViaCepResultDTO) {

	apiReturn, err := fetchAPI(url)
	if err != nil {
		fmt.Printf("Erro: %v \n", err)
		return
	}

	err = json.Unmarshal(apiReturn.Body, &dto)
	if err != nil {
		fmt.Println("Erro ao converter a resposta JSON:", err)
		return
	}

	messareReturn := MessageReturn{
		Message:  apiReturn.Message,
		Response: dto.ToString(),
	}

	ch <- messareReturn
}

func workerCdnApicep(ch chan<- interface{}, url string, dto ApicepResultDTO) {

	apiReturn, err := fetchAPI(url)
	if err != nil {
		fmt.Printf("Erro: %v \n", err)
		return
	}

	err = json.Unmarshal(apiReturn.Body, &dto)
	if err != nil {
		fmt.Println("Erro ao converter a resposta JSON:", err)
		return
	}

	messareReturn := MessageReturn{
		Message:  apiReturn.Message,
		Response: dto.ToString(),
	}

	ch <- messareReturn
}

func fetchAPI(url string) (ApiReturn, error) {
	client := http.Client{}
	startTime := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		return ApiReturn{}, err
	}
	defer resp.Body.Close()
	elapsedTime := time.Since(startTime)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ApiReturn{}, err
	}

	message := fmt.Sprintf("Tempo de resposta para %s: %s \n", url, elapsedTime)

	return ApiReturn{
		message,
		body,
	}, nil
}
