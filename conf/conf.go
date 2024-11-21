package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port    string `json:"port" env-default:"8080"`
	DSN     string `json:"dsn"`
	LogFile string `json:"log_file"`
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
