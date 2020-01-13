package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type storeConf struct {
	Path string `json:"path"`
}

// Config struct
type Config struct {
	TgToken      string    `json:"telegram_token"`
	TgEndpoint   string    `json:"telegram_endpoint"`
	DinsEndpoint string    `json:"dins_endpoint"`
	Store        storeConf `json:"store"`
}

//JSON view for config struct
func (c *Config) JSON() string {
	bytes, _ := json.Marshal(c)

	return string(bytes)

}

// FromFile read a config form json file
func FromFile(file string) *Config {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("Failed read config, check ./config.json or use flag -conf=/path, err: ", err)
	}
	var conf Config
	pErr := json.Unmarshal(bytes, &conf)
	if pErr != nil {
		log.Fatal("Failed parse config: ", pErr)
	}
	return &conf

}
