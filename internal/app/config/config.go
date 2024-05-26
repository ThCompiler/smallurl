package config

import (
	"os"
	"smallurl/pkg/logger"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Mode string

const (
	Release     Mode = "release"
	Debug       Mode = "debug"
	DebugProf   Mode = "debug+prof"
	ReleaseProf Mode = "release+prof"
)

type (
	Config struct {
		HTTP        HTTP       `yaml:"http"`
		GRPC        GRPC       `yaml:"grpc"`
		UseInMemory bool       `yaml:"use_in_memory"`
		Postgres    PG         `yaml:"postgres"`
		LoggerInfo  LoggerInfo `yaml:"logger"`
		Mode        Mode       `yaml:"mode"`
	}

	LoggerInfo struct {
		AppName           string          `yaml:"app_name"`
		Directory         string          `yaml:"directory"`
		Level             logger.LogLevel `yaml:"level"`
		UseStdAndFile     bool            `yaml:"use_std_and_file"`
		AllowShowLowLevel bool            `yaml:"allow_show_low_level"`
	}

	HTTP struct {
		Port string `yaml:"port"`
	}

	GRPC struct {
		Port     string `yaml:"port"`
		Hostname string `yaml:"hostname"`
	}

	PG struct {
		URL                string `yaml:"url"`
		MaxConnections     int    `yaml:"max_connections" default:"5"`
		MinConnections     int    `yaml:"min_connections" default:"2"`
		TTLIDleConnections uint64 `yaml:"ttl_idle_connections" default:"10"`
	}
)

func NewConfig(path string) (*Config, error) {
	cfg := &Config{}

	f, err := os.OpenFile(path, os.O_RDONLY|os.O_SYNC, 0)
	if err != nil {
		return nil, errors.Wrap(err, "can't open config file")
	}

	defer func() {
		_ = f.Close()
	}()

	if err = yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, errors.Wrap(err, "config parse error")
	}

	return cfg, nil
}
