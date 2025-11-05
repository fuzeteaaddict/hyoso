package main

import (
	"fmt"
	"log"

	"github.com/fuzeteaaddict/hyoso/internal/config"
	"github.com/fuzeteaaddict/hyoso/internal/sshd"
	"github.com/fuzeteaaddict/hyoso/internal/util"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	hostKey := util.ExpandHome(conf.Core.MasterKey)
	addr := fmt.Sprintf(":%d", conf.Core.ListenPort)

	log.Printf("[+] hyoso starting on %s ...", addr)
	log.Printf("[+] using host key: %s", hostKey)

	srv := sshd.Server{
		Addr:    addr,
		KeyPath: hostKey,
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
