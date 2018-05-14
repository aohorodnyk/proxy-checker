package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type confType struct {
	HttpbinHost    string `json:"httpbin_host"`
	Concurency     int
	Cookies        map[string]string
	FileNameSource string `json:"proxies_source"`
	FileNameResult string `json:"proxies_result"`
}

var configuration confType

func readConfiguration(filename string) {
	// Default values of the confifuration struct
	configuration = confType{
		HttpbinHost:    "https://httpbin.org/ip",
		Concurency:     100,
		FileNameSource: "./proxy.list",
		FileNameResult: "./result.list",
	}

	configData, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("Can not read", filename, err)
	}

	err = json.Unmarshal(configData, &configuration)
	if err != nil {
		log.Fatalln("Can not unmarshal", filename, err)
	}
}
