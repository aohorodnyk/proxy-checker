package core

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var sourceIP string

// Run is main function for comfortable build
func Run() {
	config := flag.String("config", "./config.json", "Path to config, can be relative")

	flag.Parse()

	readConfiguration(*config)

	proxies, err := readFileList(configuration.Source)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Loaded proxies:", len(proxies))

	sourceIP = getSourceIP()
	fmt.Println("Local IP:", sourceIP)

	var wgcp sync.WaitGroup
	var wgwf sync.WaitGroup

	in := make(chan string, configuration.Concurency)
	out := make(chan string, configuration.Concurency)

	wgwf.Add(1)
	go writeFileListChannel(out, &wgwf)

	for i := 0; i < configuration.Concurency; i++ {
		wgcp.Add(1)
		go checkProxyChannel(in, out, &wgcp)
	}

	for _, proxy := range proxies {
		in <- proxy
	}

	close(in)
	wgcp.Wait()

	close(out)
	wgwf.Wait()
}

func getSourceIP() string {
	httpClient := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	req, _ := http.NewRequest("GET", configuration.HttpbinHost, nil)

	addCookies(req)
	response, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln(configuration.HttpbinHost, err)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	var ip map[string]string
	err = json.Unmarshal(body, &ip)
	if err != nil {
		log.Fatalln(err)
	}

	proxyIP, ok := ip["origin"]
	if ok == false {
		err = proxyError{
			Message: "'origin' in JSON from IP request doesn't exists",
		}
		log.Fatalln(err)
	}

	return proxyIP
}
