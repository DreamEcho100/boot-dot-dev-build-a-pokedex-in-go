package fetch

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DreamEcho100/boot-dot-dev-build-a-pokedex-in-go/internal/cache"
)

func Get[T any](url string, cache *cache.Cache) (T, error) {
	var zero T

	var body T
	val, ok := cache.Get(url)
	if ok {
		if err := json.Unmarshal(val, &body); err != nil {
			return zero, err
		}

		return body, nil
	}

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

	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return zero, err
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return zero, err
	}
	cache.Add(url, bodyBytes)

	return body, nil
}
