package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ggrafu/sticker/services"
)

func main() {
	log.Println("starting new service")

	symbol, ok := os.LookupEnv("SYMBOL")
	if !ok {
		panic("missing env var SYMBOL")
	}
	ndays, ok := os.LookupEnv("NDAYS")
	if !ok {
		panic("missing env var NDAYS")
	}
	days, err := strconv.ParseInt(ndays, 10, 32)
	if err != nil {
		panic("incorrect value of NDAYS env var")
	}
	apiKey, ok := os.LookupEnv("APIKEY")
	if !ok {
		panic("missing env var APIKEY")
	}

	s := services.NewService(symbol, int(days), apiKey)

	http.HandleFunc("/v1/data", s.GetData)
	http.HandleFunc("/v1/ready", s.Status)
	http.ListenAndServe(":8080", nil)
}
