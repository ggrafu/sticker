package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ggrafu/sticker/services"

	rcache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
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

	redis := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	redisCache := rcache.New(&rcache.Options{
		Redis: redis,
	})

	s := services.NewService(symbol, int(days), apiKey, redisCache)

	http.HandleFunc("/v1/data", s.GetData)
	http.HandleFunc("/v1/ready", s.Status)
	http.ListenAndServe(":8080", nil)
}
