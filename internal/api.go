package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	Next     *string
	Previous *string
}

const (
	locationURL = "https://pokeapi.co/api/v2/location-area/"
	pokemonURL  = "https://pokeapi.co/api/v2/pokemon/"
)

func GetLocation(config *Config) (LocationArea, error) {
	requestURL := locationURL
	if config.Next != nil {
		requestURL = *config.Next
	}
	var object LocationArea
	err := handleGetRequest(requestURL, config, &object)
	if err != nil {
		return LocationArea{}, err
	}
	return object, nil
}

func handleGetRequest[T any](url string, config *Config, object *T) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("could not send request: %v", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not get response: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("could not read response with io: %v", err)
	}
	err = json.Unmarshal(body, &object)
	if err != nil {
		return fmt.Errorf("could not transform JSON to go struct: %v", err)
	}
	return nil
}
