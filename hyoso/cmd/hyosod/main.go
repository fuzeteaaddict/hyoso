package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fuzeteaaddict/hyoso/internal/config"
	"github.com/fuzeteaaddict/hyoso/internal/sshd"
)

func main() {
	cfgPath := filepath.Join(os.Getenv("HOME"), ".hyoso", "config.toml")

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	server := &sshd.Server{Config: cfg}
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
