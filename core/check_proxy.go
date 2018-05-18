package core

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func checkProxyChannel(in chan string, out chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for proxy := range in {

		for _, proto := range configuration.Protos {
			proxyURL := proto + "://" + proxy
			err := checkProxy(proxyURL)
			if err == nil {
				out <- proxyURL
				break
			}
		}
	}
}

func checkProxy(proxyURL string) (err error) {
	parsedProxyURL, _ := url.Parse(proxyURL)

	httpClient := &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(parsedProxyURL)},
		Timeout:   time.Duration(10 * time.Second),
	}

	req, _ := http.NewRequest("GET", configuration.HttpbinHost, nil)

	addCookies(req)

	response, err := httpClient.Do(req)
	if err != nil {
		log.Println(proxyURL, "isn't works", err)
		return err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	return checkIPResponse(body, parsedProxyURL)
}

func addCookies(req *http.Request) {
	for name, value := range configuration.Cookies {
		cookie := &http.Cookie{
			Name:    name,
			Value:   value,
			Expires: time.Now().Add(365 * 24 * time.Hour),
		}
		req.AddCookie(cookie)
	}
}

func checkIPResponse(body []byte, parsedProxyURL *url.URL) (err error) {
	var ip map[string]string
	err = json.Unmarshal(body, &ip)
	if err != nil {
		log.Println(err)
		return err
	}
	proxyIP, ok := ip["origin"]
	if ok == false {
		err = proxyError{
			Message: "'origin' in JSON from IP request doesn't exists",
		}
		log.Println(err)
		return err
	}

	if proxyIP == sourceIP {
		err = proxyError{
			Message: "Returned source IP insted of " + parsedProxyURL.Host,
		}
		log.Println(err)
		return err
	}

	if parsedProxyURL.Hostname() != proxyIP {
		err = proxyError{
			Message: proxyIP + " doesn't equal to the " + parsedProxyURL.Host,
		}
		log.Println(err)
		if configuration.ForbidMixedIP == true {
			return err
		}
	}

	return nil
}
