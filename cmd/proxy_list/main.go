package main

import (
	"fmt"

	proxyList "github.com/bcho/proxy_list"
	"github.com/bcho/proxy_list/xicidaili"
)

func main() {
	for _, proxy := range xicidaili.GetProxies() {
		ok, _ := proxyList.ValidateHTTP(proxy)
		if ok {
			fmt.Println(proxy)
		}
	}
}
