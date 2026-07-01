package fetch

import (
	"encoding/json"
	"net/http"
	"time"
)

func Get[T any](url string) (T, error) {
	var zero T

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return zero, err
	}

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return zero, err
	}
	defer res.Body.Close()

	var body T
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return zero, err
	}

	return body, nil
}
