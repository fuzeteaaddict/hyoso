package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type CoreConfig struct {
	MasterKey    string `toml:"master_key"`
	ListenPort   int    `toml:"listen_port"`
	LogDir       string `toml:"log_dir"`
	AuthMethod   string `toml:"auth_method"`
	PasswordType string `toml:"password_type"`
	PasswordFile string `toml:"password_file"`
	AuthKeyFile  string `toml:"authkey_file"`
	AuthCommand  string `toml:"auth_command"`
}

type Config struct {
	Core CoreConfig `toml:"core"`
}

func Load(path string) (*Config, error) {
	var cfg Config

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", path)
	}

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	expand := func(p string) string {
		if p == "" {
			return ""
		}
		if p[0] == '~' {
			home, _ := os.UserHomeDir()
			p = filepath.Join(home, p[1:])
		}
		return os.ExpandEnv(p)
	}

	cfg.Core.MasterKey = expand(cfg.Core.MasterKey)
	cfg.Core.LogDir = expand(cfg.Core.LogDir)
	cfg.Core.PasswordFile = expand(cfg.Core.PasswordFile)
	cfg.Core.AuthKeyFile = expand(cfg.Core.AuthKeyFile)

	if cfg.Core.ListenPort == 0 {
		cfg.Core.ListenPort = 2222
	}
	// log.Printf("loaded config: %+v", cfg.Core)
	return &cfg, nil
}
