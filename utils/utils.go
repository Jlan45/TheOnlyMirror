package utils

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
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

func GetContentReplaceReverseProxy(SourceUrl *url.URL, orgstr string, dststr string) *httputil.ReverseProxy {
	// 对内容进行替换的proxy
	proxy := httputil.NewSingleHostReverseProxy(SourceUrl)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = SourceUrl.Host
		req.Header.Set("Accept-Encoding", "deflate")
	}
	proxy.ModifyResponse = func(resp *http.Response) error {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()
		modifiedBody := string(bodyBytes)
		modifiedBody = strings.Replace(modifiedBody, orgstr, dststr, -1)
		resp.Body = io.NopCloser(strings.NewReader(modifiedBody))
		resp.ContentLength = int64(len(modifiedBody))
		return nil
	}
	return proxy
}
