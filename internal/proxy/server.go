package proxy

import (
	"fmt"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type cachedIP struct {
	ip        net.IP
	expiresAt time.Time
}

// ProxyServer represents the SOCKS5 proxy server
type ProxyServer struct {
	Port       int
	cache      map[string]cachedIP
	cacheMutex sync.RWMutex
	maxCache   int
	ttl        time.Duration
}

func NewProxyServer(port int, maxCache int, ttl time.Duration, vacuumDelay int) *ProxyServer {
	ps := &ProxyServer{
		Port:     port,
		cache:    make(map[string]cachedIP),
		maxCache: maxCache,
		ttl:      ttl * time.Second,
	}
	go ps.cacheCleaner(vacuumDelay)
	return ps
}

// Start starts the SOCKS5 server
func (ps *ProxyServer) Start() error {
	// Start TCP listener for SOCKS5 connections
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ps.Port))
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}
	defer listener.Close()

	log.Printf("SOCKS5 proxy listening on port %d", ps.Port)

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Errorf("failed to accept connection: %v", err)
			continue
		}
		log.Infoln("new client accepted")
		go ps.handleClient(client)
	}
}
