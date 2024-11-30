package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port         string       `json:"port" env-default:"8080"`
	DSN          string       `json:"dsn"`
	LogFile      string       `json:"log_file"`
	GoogleConfig GoogleConfig `json:"google_config"`
	GithubConfig GithubConfig `json:"github_config"`
}

type GoogleConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

type GithubConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

func Load(path string) (*Config, error) {
	confJSON, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	conf := &Config{}
	err = json.Unmarshal(confJSON, conf)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file: %w", err)
	}

	return conf, nil
}
