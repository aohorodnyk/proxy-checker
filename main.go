package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	proxies, err := readFile("./proxy.list")
	if err != nil {
		log.Fatalln(err)
	}

	concurrency := 300

	log.Printf("Loaded proxies: %d\n", len(proxies))

	in := make(chan string, concurrency)
	out := make(chan string, concurrency)

	var wgcp sync.WaitGroup
	var wgwf sync.WaitGroup

	wgwf.Add(1)
	go writeFile(out, &wgwf)

	for i := 0; i < concurrency; i++ {
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

func writeFile(out chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	wf := func(proxyURL string) {
		f, err := os.OpenFile("./result.list", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalln(err)
		}

		defer f.Close()

		if _, err = f.WriteString(proxyURL + "\n"); err != nil {
			log.Fatalln(err)
		}
	}

	for proxyURL := range out {
		wf(proxyURL)
	}
}

func checkProxy(in chan string, out chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for proxy := range in {

		proxyURL := "tcp://" + proxy
		parsedProxyURL, _ := url.Parse(proxyURL)
		timeout := time.Duration(10 * time.Second)
		httpClient := &http.Client{
			Transport: &http.Transport{Proxy: http.ProxyURL(parsedProxyURL)},
			Timeout:   timeout,
		}
		url := "https://httpbin.org/ip"
		req, _ := http.NewRequest("GET", url, nil)
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
			if ok != true {
				log.Println("'origin' in JSON from IP request doesn't exists")
				continue
			}

			if strings.Split(parsedProxyURL.Host, ":")[0] != proxyIP {
				log.Println(proxyIP, "Doesn't equal to the", parsedProxyURL.Host)
				continue
			}

			out <- proxyURL
		}
	}
}

func readFile(fileName string) (lines []string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// This is our buffer now
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, nil
}
