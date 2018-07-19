package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const (
	defaultAuthHTTPPort = 49152
)

// Configuration holds some configurable values
type Configuration struct {
	DeviceID         string `json:"deviceID"`
	SpotifyClientID  string `json:"spotifyClientID"`
	SpotifySecretKey string `json:"spotifySecretKey"`
	AuthHTTPPort     int    `json:"authHTTPPort"`
}

// LoadConfiguration loads the configuration from the given file.
func LoadConfiguration(filename string) *Configuration {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	config := &Configuration{
		AuthHTTPPort: defaultAuthHTTPPort,
	}
	if err = json.Unmarshal(content, &config); err != nil {
		log.Fatal(err)
	}
	return config
}
