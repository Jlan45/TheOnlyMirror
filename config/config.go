package config

import (
	"encoding/json"
	"fmt"
	"io"
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
}

var ServerConfig *Config

func (cfg *Config) GetSourceUrl(key string) *url.URL {
	tmpUrl, _ := url.Parse(cfg.Source[key])
	return tmpUrl
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
