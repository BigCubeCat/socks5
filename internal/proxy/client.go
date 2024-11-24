package proxy

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// SOCKS5 handshake
func (ps *ProxyServer) doHandshake(client net.Conn) ([]byte, error) {
	buf := make([]byte, 256)
	_, err := client.Read(buf)
	if err != nil {
		log.Errorf("invalid SOCKS5 handshake: %v", err)
		return buf, err
	}
	if buf[0] != socks5Version {
		log.Errorln("unsupported version")
	}
	// Respond with "no authentication required"
	client.Write([]byte{socks5Version, 0x00})
	return buf, nil
}

func (ps *ProxyServer) doRequestParsing(client net.Conn, buf []byte) (string, uint16, error) {
	_, err := client.Read(buf)
	if err != nil {
		log.Printf("invalid request: %v", err)
		client.Write([]byte{socks5Version, replyFailure, 0x00})
		return "", 0, err
	}
	if buf[0] != socks5Version {
		log.Errorf("unsupported version")
		client.Write([]byte{socks5Version, replyFailure, 0x00})
		return "", 0, errors.New("unsupported version")
	}
	if buf[1] != cmdConnect {
		log.Errorf("connect command expected")
		client.Write([]byte{socks5Version, replyFailure, 0x00})
		return "", 0, errors.New("connect command expected")
	}

	address, port, err := ps.parseRequest(buf)
	if err != nil {
		log.Printf("failed to parse request: %v", err)
		client.Write([]byte{socks5Version, replyFailure, 0x00})
		return "", 0, err
	}

	// резолвим доменное имя, если нужно
	if net.ParseIP(address) == nil {
		resolvedIP, err := ps.resolveDomain(address)
		log.Debugf("address=%s", address)
		if err != nil {
			log.Printf("failed to resolve domain: %v", err)
			client.Write([]byte{socks5Version, replyFailure, 0x00})
			return "", 0, err
		}
		address = resolvedIP.String()
	}
	return address, port, nil
}

// handleClient processes a SOCKS5 client connection
func (ps *ProxyServer) handleClient(client net.Conn) {
	defer client.Close()
	// Делаем Handshake
	buf, err := ps.doHandshake(client)
	if err != nil {
		return
	}
	address, port, err := ps.doRequestParsing(client, buf)
	if err != nil {
		return
	}

	// Step 3: Establish connection
	target, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", address, port), 5*time.Second)
	if err != nil {
		log.Printf("failed to connect to target: %v", err)
		client.Write([]byte{socks5Version, replyFailure, 0x00})
		return
	}
	defer target.Close()

	// Respond with success
	client.Write([]byte{socks5Version, replySuccess, 0x00, addrTypeIPv4, 0, 0, 0, 0, 0, 0})

	// Step 4: Bidirectional transfer
	ps.transferData(client, target)
}

// transferData implements bidirectional data transfer between client and target
func (ps *ProxyServer) transferData(client, target net.Conn) {
	done := make(chan struct{})

	go func() {
		defer close(done)
		io.Copy(target, client)
	}()

	io.Copy(client, target)
	<-done
}
