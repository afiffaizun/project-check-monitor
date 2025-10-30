package config

import (
	"encoding/json"
	"os"
)

type Endpoint struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Method string `json:"method"`
	Timeout int `json:"timeout"`
}

type Config struct {
	Endpoints []Endpoint `json:"endpoints"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}