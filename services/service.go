package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ggrafu/sticker/cache"
	"github.com/ggrafu/sticker/utils"
)

type Response struct {
	Values  []float32 `json:"values"`
	Average float32   `json:"average"`
}

type Service struct {
	Cache      cache.Cache
	Days       int
	APIUpdater utils.ApiUpdater
}

func NewService(symbol string, days int, apiKey string) *Service {
	return &Service{
		Cache:      *cache.NewCache(),
		Days:       days,
		APIUpdater: utils.NewAVUpdater(apiKey, symbol),
	}
}

// function GetData is http handler to use as a main data provider
func (s *Service) GetData(w http.ResponseWriter, r *http.Request) {

	e := json.NewEncoder(w)

	if s.Cache.IsOutdated() {
		log.Println("cache is empty or outdated. Updating the cache...")
		ts, err := s.APIUpdater.FetchAPI()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// update cache
		err = s.Cache.Update(ts)
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
