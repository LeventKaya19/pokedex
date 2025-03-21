package mapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type AreaResults struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetMap(url string) (AreaResults, error) {
	res, err := http.Get(url)
	if err != nil {
		return AreaResults{}, errors.New("Fatal error")
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return AreaResults{}, fmt.Errorf("Bad response: %s", res.Status)
	}

	var result AreaResults
	err = json.Unmarshal(body, &result)
	if err != nil {
		return AreaResults{}, fmt.Errorf("Could not unmarshal")
	}

	return result, nil
}
