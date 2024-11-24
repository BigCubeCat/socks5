package main

import (
	"fmt"
	"os"
	"socks/internal/proxy"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <port>\n", os.Args[0])
		os.Exit(1)
	}
	log.SetLevel(log.TraceLevel)

	port, err := strconv.Atoi(os.Args[1])
	if err != nil || port <= 0 || port > 65535 {
		log.Fatalf("Invalid port: %s", os.Args[1])
	}

	var d time.Duration
	d = 1000000000
	server := proxy.NewProxyServer(port, 100, d)
	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
