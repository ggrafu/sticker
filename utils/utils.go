package utils

import (
	"bufio"
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	API_URL         string = "https://www.alphavantage.co/query"
	API_FUNCTION    string = "TIME_SERIES_DAILY_ADJUSTED"
	API_OUTPUT_SIZE string = "full"
)

type ApiUpdater interface {
	FetchAPI() (*APIData, error)
}

type AVUpdater struct {
	APIKey string
	Symbol string
}

func NewAVUpdater(apiKey string, symbol string) *AVUpdater {
	return &AVUpdater{
		APIKey: apiKey,
		Symbol: symbol,
	}
}

func (a *AVUpdater) FetchAPI() (*APIData, error) {
	req, err := http.NewRequest("GET", API_URL, nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Add("apikey", a.APIKey)
	q.Add("function", API_FUNCTION)
	q.Add("symbol", a.Symbol)
	q.Add("outputsize", API_OUTPUT_SIZE)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ts APIData
	d := json.NewDecoder(bufio.NewReader(resp.Body))
	if err := d.Decode(&ts); err != nil {
		return nil, err
	}
	return &ts, nil
}
