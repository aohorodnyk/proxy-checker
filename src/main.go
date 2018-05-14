package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var localIP string

func main() {
	readConfiguration("./config.json")

	proxies, err := readFileList(configuration.FileNameSource)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Loaded proxies:", len(proxies))

	var wgcp sync.WaitGroup
	var wgwf sync.WaitGroup

	in := make(chan string, configuration.Concurency)
	out := make(chan string, configuration.Concurency)

	wgwf.Add(1)
	go writeFileListChannel(out, &wgwf)

	for i := 0; i < configuration.Concurency; i++ {
		wgcp.Add(1)
		go checkProxy(in, out, &wgcp)
	}

	for _, proxy := range proxies {
		in <- proxy
	}

	close(in)
	wgcp.Wait()

	close(out)
	wgwf.Wait()
}

func checkProxy(in chan string, out chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for proxy := range in {

		proxyURL := "tcp://" + proxy
		parsedProxyURL, _ := url.Parse(proxyURL)

		httpClient := &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(parsedProxyURL)},
			Timeout:   time.Duration(10 * time.Second),
		}

		req, _ := http.NewRequest("GET", configuration.HttpbinHost, nil)

		for name, value := range configuration.Cookies {
			cookie := &http.Cookie{
				Name:    name,
				Value:   value,
				Expires: time.Now().Add(365 * 24 * time.Hour),
				Domain:  "httpbin.d.ohorodnyk.name",
			}
			req.AddCookie(cookie)
		}

		response, err := httpClient.Do(req)
		if err != nil {
			log.Println(proxy, " isn't works ", err)
		} else {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)

			var ip map[string]string
			err := json.Unmarshal(body, &ip)
			if err != nil {
				log.Println(err)
				continue
			}
			proxyIP, ok := ip["origin"]
			if ok == false {
				log.Println("'origin' in JSON from IP request doesn't exists")
				continue
			}

			if parsedProxyURL.Hostname() != proxyIP {
				log.Println(proxyIP, "Doesn't equal to the", parsedProxyURL.Host)
				continue
			}

			out <- proxyURL
		}
	}
}
