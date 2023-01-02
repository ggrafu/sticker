package utils

type APIData struct {
	Metadata   map[string]interface{} `json:"Meta Data"`
	TimeSeries map[string]Record      `json:"Time Series (Daily)"`
}

type Record struct {
	Close string `json:"4. close"`
}
