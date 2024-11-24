package proxy

import (
	"context"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// resolveDomain resolves a domain name to an IP address using the standard resolver
func (ps *ProxyServer) resolveDomain(domain string) (net.IP, error) {
	ps.cacheMutex.RLock()
	ipEntry, found := ps.cache[domain]
	ps.cacheMutex.RUnlock()

	if found && time.Now().Before(ipEntry.expiresAt) {
		log.Debugf("Cache hit for domain: %s -> %s", domain, ipEntry.ip)
		return ipEntry.ip, nil
	}
	// Если не найдено или запись устарела, выполняем запрос DNS

	ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", domain)
	if err != nil {
		log.Errorf("failed to resolve domain: %s", err.Error())
		return nil, fmt.Errorf("failed to resolve domain: %s", err.Error())
	}

	var resolvedIP net.IP
	for _, ip := range ips {
		if ip.To4() != nil {
			log.Debugf("Domain resolved: %s=%s", domain, ip.String())
			resolvedIP = ip
			break
		}
	}
	if resolvedIP == nil {
		log.Errorf("No IPv4 address found for domain: %s", domain)
		return nil, fmt.Errorf("no IPv4 address found for domain: %s", domain)
	}
	// Кэшируем результат
	ps.cacheMutex.Lock()
	ps.cache[domain] = cachedIP{
		ip:        resolvedIP,
		expiresAt: time.Now().Add(ps.ttl),
	}
	ps.cacheMutex.Unlock()

	log.Debugf("Domain resolved and cached: %s -> %s", domain, resolvedIP)
	return resolvedIP, nil
}
