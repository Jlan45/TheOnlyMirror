package utils

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func GetSimpleReverseProxy(SourceUrl *url.URL) *httputil.ReverseProxy {
	// 最简单的proxy，只进行host替换，并且去除gzip
	proxy := httputil.NewSingleHostReverseProxy(SourceUrl)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = SourceUrl.Host
		req.Header.Set("Accept-Encoding", "deflate")
	}
	return proxy
}
