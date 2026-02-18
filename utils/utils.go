package utils

import (
	"TheOnlyMirror/config"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// maxReplaceBodySize 是内容替换的最大响应体大小 (10MB)。
// 超过此大小的响应将直接透传，不做替换，防止 OOM。
const maxReplaceBodySize = 10 << 20

// buildTransport 根据 sourceKey 创建带代理配置的 Transport。
func buildTransport(sourceKey string) *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	if proxyFunc := config.ServerConfig.GetProxyForSource(sourceKey); proxyFunc != nil {
		t.Proxy = proxyFunc
		log.Printf("[%s] using proxy: %s", sourceKey, config.ServerConfig.Proxy)
	}
	return t
}

func GetSimpleReverseProxy(SourceUrl *url.URL, sourceKey string) *httputil.ReverseProxy {
	// 最简单的proxy，只进行host替换
	log.Printf("[%s] creating simple reverse proxy -> %s", sourceKey, SourceUrl.String())
	proxy := httputil.NewSingleHostReverseProxy(SourceUrl)
	proxy.Transport = buildTransport(sourceKey)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = SourceUrl.Host
	}
	return proxy
}

func GetContentReplaceReverseProxy(SourceUrl *url.URL, orgstr string, dststr string, sourceKey string) *httputil.ReverseProxy {
	// 对内容进行替换的proxy
	log.Printf("[%s] creating content-replace reverse proxy -> %s (replace: %s => %s)", sourceKey, SourceUrl.String(), orgstr, dststr)
	proxy := httputil.NewSingleHostReverseProxy(SourceUrl)
	proxy.Transport = buildTransport(sourceKey)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = SourceUrl.Host
		// 使用 identity 确保上游返回未压缩内容，避免在压缩数据上做字符串替换导致损坏
		req.Header.Set("Accept-Encoding", "identity")
	}
	proxy.ModifyResponse = func(resp *http.Response) error {
		// 如果上游仍返回了压缩内容，跳过替换直接透传
		if enc := resp.Header.Get("Content-Encoding"); enc != "" && enc != "identity" {
			log.Printf("warn: upstream returned Content-Encoding=%s, skipping replacement for %s", enc, resp.Request.URL)
			return nil
		}

		// 对超大响应跳过替换，防止 OOM
		if resp.ContentLength > maxReplaceBodySize {
			log.Printf("warn: response too large (%d bytes), skipping replacement for %s", resp.ContentLength, resp.Request.URL)
			return nil
		}

		// 使用 LimitReader 兜底，即使 ContentLength 未设置或不可信也不会超限
		limitedReader := io.LimitReader(resp.Body, maxReplaceBodySize+1)
		bodyBytes, err := io.ReadAll(limitedReader)
		if err != nil {
			return err
		}
		resp.Body.Close()

		if int64(len(bodyBytes)) > maxReplaceBodySize {
			// 读取已超限，直接返回原始内容
			log.Printf("warn: response exceeded %d bytes limit, skipping replacement", maxReplaceBodySize)
			resp.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
			return nil
		}

		modifiedBody := strings.ReplaceAll(string(bodyBytes), orgstr, dststr)
		resp.Body = io.NopCloser(strings.NewReader(modifiedBody))
		resp.ContentLength = int64(len(modifiedBody))
		resp.Header.Set("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
		// 替换后内容已经是完整的非 chunked 数据，必须移除 Transfer-Encoding
		resp.Header.Del("Transfer-Encoding")
		return nil
	}
	return proxy
}
