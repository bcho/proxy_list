package xicidaili

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

const proxiesListUrl = "http://www.xicidaili.com"

// GetProxies returns list of proxy ip from xicidaili.com .
func GetProxies() []*url.URL {
	req, err := http.NewRequest("GET", proxiesListUrl, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "curl/7.49.1")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil
	}

	var proxies []*url.URL

	doc.Find("#ip_list tr").Each(func(_ int, s *goquery.Selection) {
		ip := s.Find("td:nth-child(2)").Text()
		port := s.Find("td:nth-child(3)").Text()
		if ip == "" || port == "" {
			return
		}

		if u, err := url.Parse(fmt.Sprintf("%s:%s", ip, port)); err == nil {
			proxies = append(proxies, u)
		}

	})

	return proxies
}
