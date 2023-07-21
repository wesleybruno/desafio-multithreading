package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

type MessageReturn struct {
	Message string
}

func main() {

	cep := flag.String("cep", "", "URL da primeira API")
	flag.Parse()

	urlAPI1 := "https://cdn.apicep.com/file/apicep/" + *cep + ".json"
	urlAPI2 := "http://viacep.com.br/ws/" + *cep + "/json/"
	timeout := 1 * time.Second

	ch := make(chan string)
	ch2 := make(chan string)

	go worker(ch, urlAPI1)
	go worker(ch2, urlAPI2)

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

func worker(ch chan<- string, url string) {

	message, err := fetchAPI(url)
	if err != nil {
		return
	}

	ch <- message.Message

}

func fetchAPI(url string) (MessageReturn, error) {
	client := http.Client{}
	startTime := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		return MessageReturn{}, err
	}
	defer resp.Body.Close()
	elapsedTime := time.Since(startTime)
	message := fmt.Sprintf("Tempo de resposta para %s: %s \n", url, elapsedTime)
	return MessageReturn{
		message,
	}, nil
}
