package sshd

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

type Server struct {
	Addr    string
	KeyPath string
}

func (s *Server) Start() error {
	privateBytes, err := os.ReadFile(s.KeyPath)
	if err != nil {
		return fmt.Errorf("read host key: %w", err)
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return fmt.Errorf("parse host key: %w", err)
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true, // todo: actual auth lmao
	}

	config.AddHostKey(private)

	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	log.Printf("Hyoso SSH listening on %s", s.Addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go handleConn(conn, config)
	}
}

func handleConn(nConn net.Conn, config *ssh.ServerConfig) {
	sshConn, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		log.Printf("failed handshake: %v", err)
		return
	}
	log.Printf("[+] new ssh conn from: %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())

	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unsupported channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("channel accept: %v", err)
			continue
		}
		go func(in <-chan *ssh.Request) {
			for req := range in {
				if req.Type == "shell" || req.Type == "exec" {
					req.Reply(true, nil)
					io.WriteString(channel, "~ welcome to hyoso ~\n")
					io.WriteString(channel, "work in progress\n")
					channel.Close()
				} else {
					req.Reply(false, nil)
				}
			}
		}(requests)
	}
}
