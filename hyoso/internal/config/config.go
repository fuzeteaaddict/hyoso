package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Core struct {
	MasterKey  string `toml:"master_key"`
	ListenPort int    `toml:"listen_port"`
	LogDir     string `toml:"log_dir"`
}

type Target struct {
	Name     string   `toml:"name"`
	Host     string   `toml:"host"`
	User     string   `toml:"user"`
	Port     int      `toml:"port"`
	KeyPath  string   `toml:"key_path"`
	Record   string   `toml:"record"`
	Tags     []string `toml:"tags"`
}

type Config struct {
	Core    Core      `toml:"core"`
	Targets []Target  `toml:"target"`
}

func LoadConfig() (*Config, error) {
	path := filepath.Join(os.Getenv("HOME"), ".hyoso", "config.toml")
	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("config not found: %w", err)
	}

	var conf Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if conf.Core.ListenPort == 0 {
		conf.Core.ListenPort = 2222
	}
	return &conf, nil
}
