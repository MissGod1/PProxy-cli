package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Server   string   `json:"server"`
	Port     uint16     `json:"port"`
	Password string   `json:"password"`
	Method   string   `json:"method"`
	Apps     []string `json:"apps"`
}

func NewConfig(path string) *Config {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("Open Configure File Err:", err)
	}
	conf := Config{}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalln("Json Configure Err:", err)
	}
	return &conf
}
