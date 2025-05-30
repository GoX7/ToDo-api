package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Cookie cookie `yaml:"cookie"`
	Server server `yaml:"server"`
	Logger logger `yaml:"logger"`
}

type cookie struct {
	Key string `yaml:"key"`
}

type server struct {
	Addr string        `json:"addr"`
	Wto  time.Duration `json:"wto"`
	Rto  time.Duration `json:"rto"`
}

type logger struct {
	Level  string `yaml:"level"`
	Server string `yaml:"server_path"`
	Sqlite string `yaml:"sqlite_path"`
	MW     string `yaml:"mw_path"`
}

var (
	cfg = Config{}
)

func Load() (*Config, error) {
	err := cleanenv.ReadConfig("config/config.yaml", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
