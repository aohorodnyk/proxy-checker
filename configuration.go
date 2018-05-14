package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type confType struct {
	HttpbinHost string `json:"httpbin_host"`
	Concurency  int
	Cookies     map[string]string
}

var configuration confType

func readConfiguration(filename string) {
	// Default values of the confifuration struct
	configuration = confType{
		HttpbinHost: "https://httpbin.org/ip",
		Concurency:  100,
	}

	configData, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("Can not read config.json", err)
	}

	err = json.Unmarshal(configData, &configuration)
	if err != nil {
		log.Fatalln("Can not unmarshal config.json", err)
	}
}
