package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type confType struct {
	HttpbinHost   string `json:"httpbin_host"`
	Concurency    int
	Cookies       map[string]string
	Source        string
	Result        string
	Protos        []string
	ForbidMixedIP bool `json:"forbid_mixed_ip"`
}

var configuration confType

func readConfiguration(filename string) {
	// Default values of the confifuration struct
	configuration = confType{
		HttpbinHost:   "https://httpbin.org/ip",
		Concurency:    100,
		Source:        "./source.list",
		Result:        "./result.list",
		Protos:        []string{"http", "tcp"},
		ForbidMixedIP: false,
	}

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println(filename + " not found")
		fmt.Println("Using the default configuration")
		return
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
