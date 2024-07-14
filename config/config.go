package config

import (
	"encoding/json"
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
	Tls          bool              `json:"tls"`
	HostControll bool              `json:"hostControll"` //是否启用访问来源host控制
	HostList     []string          `json:"hostList"`     //允许访问的host列表
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
	return nil
}
