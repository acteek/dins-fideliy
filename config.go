package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type StoreConf struct {
	Path string `json:"path"`
}

type Config struct {
	TgToken      string    `json:"telegram_token"`
	TgEndpoint   string    `json:"telegram_endpoint"`
	DinsEndpoint string    `json:"dins_endpoint"`
	Store        StoreConf `json:"store"`
}

func (c *Config) Json() string {
	bytes, _ := json.Marshal(c)

	return string(bytes)

}

func FromFile() *Config {
	bytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Failed Read ./conf.json", err)
	}
	var conf Config
	pErr := json.Unmarshal(bytes, &conf)
	if pErr != nil {
		log.Fatal("Failed Parse ./conf.json", pErr)
	}
	return &conf

}
