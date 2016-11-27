package proxy_list

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

const validateUrl = "http://httpbin.org/ip"

var (
	validateTimeout = time.Duration(10 * time.Second)

	errInvalidProxy = errors.New("invalid proxy")
)

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

	resp, err := client.Get(validateUrl)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 {
		return false, errInvalidProxy
	}

	return true, nil
}
