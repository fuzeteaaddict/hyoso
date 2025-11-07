package sshd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fuzeteaaddict/hyoso/internal/config"
	gliderssh "github.com/gliderlabs/ssh"
)

type Server struct {
	Config *config.Config
}

func (s *Server) Start() error {
	c := s.Config.Core

	addr := fmt.Sprintf(":%d", c.ListenPort)
	hostKey := c.MasterKey
	authMethod := strings.ToLower(c.AuthMethod)

	if _, err := os.Stat(hostKey); err != nil {
		return fmt.Errorf("missing master key: %v", err)
	}

	var opts []gliderssh.Option

	opts = append(opts, gliderssh.HostKeyFile(hostKey))

	switch authMethod {
	case "pubkey":
		opts = append(opts, gliderssh.PublicKeyAuth(s.pubkeyAuth))
	case "password":
		opts = append(opts, gliderssh.PasswordAuth(s.passwordAuth))
	case "custom":
		opts = append(opts, gliderssh.PasswordAuth(s.customAuth))
	default:
		opts = append(opts, gliderssh.PublicKeyAuth(s.pubkeyAuth))
	}

	gliderssh.Handle(func(sess gliderssh.Session) {
		user := sess.User()
		cmd := sess.Command()

		fmt.Fprintf(sess, "logged in as %s\n", user)

		if len(cmd) > 0 {
			fmt.Fprintf(sess, "running command: %v\n", cmd)
			return
		}

		io.WriteString(sess, "~ welcome to hyoso ~\n")
		io.WriteString(sess, "exiting\n")
	})

	log.Printf("[+] Starting Hyoso daemon on %s (auth=%s)", addr, authMethod)
	return gliderssh.ListenAndServe(addr, nil, opts...)
}

func (s *Server) pubkeyAuth(ctx gliderssh.Context, key gliderssh.PublicKey) bool {
	path := s.Config.Core.AuthKeyFile
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("[auth] failed to read authkey_file: %v", err)
		return false
	}

	for len(data) > 0 {
		pub, _, _, rest, err := gliderssh.ParseAuthorizedKey(data)
		if err != nil {
			break
		}
		if gliderssh.KeysEqual(pub, key) {
			log.Printf("[auth] pubkey accepted for %s", ctx.User())
			return true
		}
		data = rest
	}

	log.Printf("[auth] pubkey rejected for %s", ctx.User())
	return false
}

func (s *Server) passwordAuth(ctx gliderssh.Context, password string) bool {
	cfg := s.Config.Core
	data, err := os.ReadFile(cfg.PasswordFile)
	if err != nil {
		log.Printf("[auth] password file read error: %v", err)
		return false
	}
	expected := strings.TrimSpace(string(data))
	// log.Printf("[auth] expected password value: %v", expected) // useful for debugging
	switch strings.ToLower(cfg.PasswordType) {
	case "plaintext", "":
		return password == expected
	case "sha256":
		h := sha256.Sum256([]byte(password))
		return hex.EncodeToString(h[:]) == expected
	// todo: add sha512, others
	default:
		log.Printf("[auth] unsupported hash type %q", cfg.PasswordType)
		return false
	}
}

func (s *Server) customAuth(ctx gliderssh.Context, password string) bool {
	cmdLine := s.Config.Core.AuthCommand
	if cmdLine == "" {
		log.Printf("[auth] custom auth command not set")
		return false
	}

	args := strings.Fields(cmdLine)
	cmd := exec.Command(args[0], args[1:]...)

	// pass as variables, not env
	// important because we don't want to expose passwords to entire system
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("HYOSO_USER=%s", ctx.User()),
		fmt.Sprintf("HYOSO_PASSWORD=%s", password),
	)

	if err := cmd.Run(); err != nil {
		log.Printf("[auth] custom command failed: %v", err)
		return false
	}

	return true
}
