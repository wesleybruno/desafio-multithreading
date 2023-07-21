package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {

	cep := flag.String("cep", "", "URL da primeira API")
	flag.Parse()

	urlAPI1 := "https://cdn.apicep.com/file/apicep/" + *cep + ".json"
	urlAPI2 := "http://viacep.com.br/ws/" + *cep + "/json/"
	timeout := 1 * time.Second

	ch := make(chan string)
	ch2 := make(chan string)

	go fetchAPI(urlAPI1, timeout, ch)
	go fetchAPI(urlAPI2, timeout, ch2)

	select {
	case result := <-ch:
		fmt.Println(result)
	case result := <-ch2:
		fmt.Println(result)
	case <-time.After(timeout):
		fmt.Println("Erro: timeout ao aguardar a resposta.")
		return
	}

}

func fetchAPI(url string, timeout time.Duration, ch chan<- string) {
	client := http.Client{
		Timeout: timeout,
	}
	startTime := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("Erro na requisição para %s: %s \n", url, err.Error())
		return
	}
	defer resp.Body.Close()
	elapsedTime := time.Since(startTime)
	ch <- fmt.Sprintf("Tempo de resposta para %s: %s \n", url, elapsedTime)
}
