package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	Next      *string
	Previous  *string
	Cache     *Cache
	Arguments []string
	Pokemons  map[string]Pokemon
	History   []string
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

func GetLocationDetail(config *Config, query string) (LocationAreaDetail, error) {
	requestURL := locationURL + query
	var object LocationAreaDetail
	err := handleGetRequest(requestURL, config, &object)
	if err != nil {
		return LocationAreaDetail{}, err
	}
	return object, nil
}

func GetPokemon(config *Config, query string) (Pokemon, error) {
	requestURL := pokemonURL + query
	var pokemon Pokemon
	err := handleGetRequest(requestURL, config, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}
	return pokemon, nil
}

func handleGetRequest[T any](url string, config *Config, object *T) error {
	if data, ok := config.Cache.Get(url); ok {
		err := json.Unmarshal(data, &object)
		if err != nil {
			return fmt.Errorf("could not transform cached JSON to go struct: %v", err)
		}
		return nil
	}

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

	config.Cache.Add(url, body)

	err = json.Unmarshal(body, &object)
	if err != nil {
		return fmt.Errorf("could not transform JSON to go struct: %v", err)
	}
	return nil
}
