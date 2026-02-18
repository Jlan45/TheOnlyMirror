package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Domain       string            `json:"domain"`
	MirrorList   []string          `json:"mirrorList"`
	Source       map[string]string `json:"sources"`
	Port         int               `json:"port"`
	TlsPort      int               `json:"tlsport"`
	Tls          bool              `json:"tls"`
	HostControll bool              `json:"hostControll"` //是否启用访问来源host控制
	HostList     []string          `json:"hostList"`     //允许访问的host列表
	CertFile     string            `json:"certFile"`
	KeyFile      string            `json:"keyFile"`
	Proxy        string            `json:"proxy"` // 全局上层代理地址 (支持 http/https/socks5)
}

var ServerConfig *Config

func (cfg *Config) GetSourceUrl(key string) *url.URL {
	tmpUrl, err := url.Parse(cfg.Source[key])
	if err != nil {
		log.Printf("error: invalid source URL for %s: %v", key, err)
		return &url.URL{}
	}
	return tmpUrl
}

// GetProxyForSource 返回指定源的代理配置函数，供 http.Transport 使用。
// 使用全局 proxy 配置，没有则返回 nil（直连）。
func (cfg *Config) GetProxyForSource(sourceKey string) func(*http.Request) (*url.URL, error) {
	if cfg.Proxy == "" {
		return nil
	}
	proxyURL, err := url.Parse(cfg.Proxy)
	if err != nil {
		log.Printf("error: invalid proxy URL: %v", err)
		return nil
	}
	return http.ProxyURL(proxyURL)
}
func (cfg *Config) CheckHost(host string) bool {
	for _, v := range cfg.HostList {
		if strings.HasPrefix(host, v) {
			return true
		}
	}
	return false
}
func (cfg *Config) GetMyUrl() *url.URL {
	tmpUrl := &url.URL{
		Scheme: "http",
		Host:   cfg.Domain,
		Path:   "",
	}
	if cfg.Tls {
		tmpUrl.Scheme = "https"
	}
	return tmpUrl
}
func LoadConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return err

	}
	err = json.Unmarshal(bytes, &ServerConfig)
	if err != nil {
		return err
	}
	// Validate all source URLs during config loading
	for key, urlStr := range ServerConfig.Source {
		if _, err := url.Parse(urlStr); err != nil {
			return fmt.Errorf("invalid URL in config for source %s: %w", key, err)
		}
	}
	return nil
}
