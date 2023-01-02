package services

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ggrafu/sticker/utils"
)

type MockUpdater struct{}

func (*MockUpdater) FetchAPI() (*utils.APIData, error) {
	return &utils.APIData{
		Metadata: nil,
		TimeSeries: map[string]utils.Record{
			"2022-12-30": {Close: "13"},
			"2022-12-29": {Close: "14"},
			"2022-12-28": {Close: "15"},
		}}, nil
}

func TestGetData(t *testing.T) {

	s := NewService("", 3, "")
	s.APIUpdater = &MockUpdater{}

	r := httptest.NewRequest(http.MethodGet, "/v1/data", nil)
	w := httptest.NewRecorder()
	s.GetData(w, r)

	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected response from GetData handler: %v", err)
	}

	exp := "{\"values\":[13,14,15],\"average\":14}\n"
	if string(data) != exp {
		t.Errorf("expected %s, but actual data: %s", exp, data)
	}
}
