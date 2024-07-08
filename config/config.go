package config

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
)

type Config struct {
	Domain     string            `json:"domain"`
	MirrorList []string          `json:"mirrorList"`
	Source     map[string]string `json:"sources"`
	Tls        bool              `json:"tls"`
}

var ServerConfig *Config

func (cfg *Config) GetSourceUrl(key string) *url.URL {
	tmpUrl, _ := url.Parse(cfg.Source[key])
	return tmpUrl
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
