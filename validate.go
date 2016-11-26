package proxy_list

import (
	"net/http"
	"net/url"
	"time"
)

const validateUrl = "http://httpbin.org/ip"

var validateTimeout = time.Duration(10 * time.Second)

// ValidateHTTP validates a http proxy.
func ValidateHTTP(proxy *url.URL) (bool, error) {
	p := &url.URL{
		Scheme: "http",
		Host:   proxy.Host,
		User:   proxy.User,
		Path:   proxy.Path,
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(p),
		},
		Timeout: validateTimeout,
	}

	_, err := client.Get(validateUrl)
	return err == nil, err
}
