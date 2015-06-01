package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Endpoint        string
	BrowserEndpoint string
	SSO             string
	Dockerhub       string
	Redis_addr      string
	DB_addr         string
	App_id          string
	App_key         string
}

func ReadConfigure(name string) (Configuration, error) {
	file, err := os.Open(name)
	configuration := Configuration{}
	if err != nil {
		return configuration, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return configuration, err
	}
	return configuration, nil
}
