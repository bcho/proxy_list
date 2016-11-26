package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	proxyList "github.com/bcho/proxy_list"
	"github.com/bcho/proxy_list/xicidaili"
	"github.com/elazarl/goproxy"
)

var refreshDuration = time.Duration(15 * time.Minute)

func main() {
	go refreshProxies(refreshDuration)

	startProxyServer(os.Getenv("PROXY_LIST_BIND"))
}

func startProxyServer(bind string) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			p := GetProxy()
			if p == nil {
				log.Printf("no proxy")
				return r, nil
			}

			proxy.Tr = &http.Transport{Proxy: http.ProxyURL(p)}

			return r, nil
		},
	)

	log.Fatal(http.ListenAndServe(bind, proxy))
}

func refreshProxies(refreshDuration time.Duration) {
	refresh := func() {
		for _, proxy := range xicidaili.GetProxies() {
			ok, err := proxyList.ValidateHTTP(proxy)
			if ok {
				log.Printf("proxy %s is usable", proxy)
				AddProxy(proxy)
			} else {
				InvalidateProxy(proxy)
				log.Printf("proxy %s is invalid: %s", proxy, err)
			}
		}
	}

	refresh()

	ticker := time.NewTicker(refreshDuration)
	defer ticker.Stop()
	for range ticker.C {
		refresh()
	}
}

type ProxyStats struct {
	url        *url.URL
	updatedAt  time.Time
	lastUsedAt time.Time
	usedTimes  int
}

type Proxies struct {
	proxies map[string]*ProxyStats
	lock    *sync.Mutex
}

var defaultProxies = &Proxies{
	proxies: make(map[string]*ProxyStats),
	lock:    &sync.Mutex{},
}

func (p *Proxies) Add(proxy *url.URL) {
	key := proxyKey(proxy)
	p.lock.Lock()
	defer p.lock.Unlock()

	if stat, present := p.proxies[key]; present {
		stat.updatedAt = time.Now()
	} else {
		p.proxies[key] = &ProxyStats{
			url: &url.URL{
				Scheme: "http",
				Host:   proxy.Host,
				User:   proxy.User,
				Path:   proxy.Path,
			},
			updatedAt: time.Now(),
			usedTimes: 0,
		}
	}
}

func (p *Proxies) Invalidate(proxy *url.URL) {
	key := proxyKey(proxy)
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(p.proxies, key)
}

func (p *Proxies) Get() *url.URL {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, stat := range p.proxies {
		stat.usedTimes = stat.usedTimes + 1
		stat.lastUsedAt = time.Now()
		return stat.url
	}

	return nil
}

func proxyKey(proxy *url.URL) string { return proxy.String() }

func AddProxy(proxy *url.URL) {
	defaultProxies.Add(proxy)
}

func InvalidateProxy(proxy *url.URL) {
	defaultProxies.Invalidate(proxy)
}

func GetProxy() *url.URL {
	return defaultProxies.Get()
}
