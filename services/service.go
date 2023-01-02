package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ggrafu/sticker/cache"
	"github.com/ggrafu/sticker/utils"

	rcache "github.com/go-redis/cache/v8"
)

type Response struct {
	Values  []float32 `json:"values"`
	Average float32   `json:"average"`
}

type Service struct {
	Cache      cache.Cache
	RCache     *rcache.Cache
	Days       int
	Symbol     string
	APIUpdater utils.ApiUpdater
}

const REDIS_CACHE_TTL = time.Hour

func NewService(symbol string, days int, apiKey string, rcache *rcache.Cache) *Service {
	return &Service{
		Cache:      *cache.NewCache(),
		RCache:     rcache,
		Days:       days,
		Symbol:     symbol,
		APIUpdater: utils.NewAVUpdater(apiKey, symbol),
	}
}

// function GetData is http handler to use as a main data provider
func (s *Service) GetData(w http.ResponseWriter, r *http.Request) {

	e := json.NewEncoder(w)

	// try to hit local cache
	if s.Cache.IsOutdated() {
		log.Println("local cache is empty or outdated. Updating the cache...")

		ctx := context.TODO()
		var cached *utils.APIData

		// local cache is invalid, try to hit redis cache if it's enabled
		if s.RCache != nil {
			err := s.RCache.Get(ctx, s.Symbol, &cached)
			// if redis cache is empty - fetching data from source API and updating the redis cache
			if err != nil {
				log.Println("redis cache is empty. Fetching data from source API...")
				cached, err = s.APIUpdater.FetchAPI()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				s.RCache.Set(&rcache.Item{
					Ctx:   ctx,
					Key:   s.Symbol,
					Value: *cached,
					TTL:   REDIS_CACHE_TTL,
				})
			}
		} else {
			// if redis cache is not enabled - just fetch source api
			c, err := s.APIUpdater.FetchAPI()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			cached = c
		}

		// update cache
		err := s.Cache.Update(cached)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// generate response slice and calculate average
	values := s.Cache.GetLastElements(s.Days)
	if len(values) == 0 {
		log.Println("unexpected response from the backend service. Please check the throttling.")
		http.Error(w, "No data available", http.StatusServiceUnavailable)
		return
	}
	size := len(values)
	var sum float32
	for i := 0; i < size; i++ {
		sum += values[i]
	}
	avg := sum / float32(size)
	response := Response{
		Values:  values,
		Average: avg,
	}
	if err := e.Encode(response); err != nil {
		log.Fatalln("failed to encode data")
	}
}

func (*Service) Status(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}
