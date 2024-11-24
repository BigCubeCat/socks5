package main

import (
	"fmt"
	"os"
	"socks/internal/proxy"
	"time"

	"github.com/akamensky/argparse"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Create new parser object
	parser := argparse.NewParser("socks", "Stupid SOCKS5 proxy")
	var port *int = parser.Int("p", "port", &argparse.Options{Required: true, Help: "port"})
	var cacheSize *int = parser.Int("c", "cache-size", &argparse.Options{Required: false, Default: 1024, Help: "maximum count of cached domains"})
	var ttl *int = parser.Int("t", "ttl", &argparse.Options{Required: false, Default: 20, Help: "time to live for cache entry, in seconds"})
	var cleanerDuration *int = parser.Int("v", "vacuum", &argparse.Options{Required: false, Default: 1, Help: "cache cleaner duration (in minutes). Used to control cache size"})
	var logLevel *string = parser.String("l", "log-level", &argparse.Options{Required: false, Default: "error", Help: "set log level; default = error "})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	loggingLevel, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("cant invalid logging level: %s", err.Error())
		return
	}
	log.SetLevel(loggingLevel)

	server := proxy.NewProxyServer(
		*port,
		*cacheSize,
		time.Duration(*ttl),
		*cleanerDuration,
	)
	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
