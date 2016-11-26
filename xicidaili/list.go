package xicidaili

import (
	"bufio"
	"net/http"
	"net/url"
)

const proxiesListUrl = "http://api.xicidaili.com/free2016.txt"

// GetProxies returns list of proxy ip from xicidaili.com .
func GetProxies() []*url.URL {
	resp, err := http.Get(proxiesListUrl)
	if err != nil {
		return nil
	}

	if resp.StatusCode != 200 {
		return nil
	}

	scanner := bufio.NewScanner(resp.Body)
	defer resp.Body.Close()

	var proxies []*url.URL
	for scanner.Scan() {
		proxyUrl, err := url.Parse(scanner.Text())
		if err == nil {
			proxies = append(proxies, proxyUrl)
		}
	}

	return proxies
}
