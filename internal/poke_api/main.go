package poke_api

type PokeAPILocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokeAPILocationAreaResultsList struct {
	Count    int                   `json:"count"`
	Next     *string               `json:"next"` // Use pointer to handle null values
	Previous *string               `json:"previous"`
	Results  []PokeAPILocationArea `json:"results"`
}

const PokeAPILocationAreaBaseURL = "https://pokeapi.co/api/v2/location-area/?limit=20&offset=0"
