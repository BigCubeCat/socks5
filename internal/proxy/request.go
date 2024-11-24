package proxy

import (
	"encoding/binary"
	"errors"
	"net"

	log "github.com/sirupsen/logrus"
)

// предназначен для обработки данных запроса SOCKS5.
// извлекает адрес и порт целевого соединения
// По сути интерпретирует информацию из клиентского запроса и подготавливает
// данные для установления соединения с сервером
func (ps *ProxyServer) parseRequest(buf []byte) (string, uint16, error) {
	var (
		address string
		port    uint16
	)
	// извлечение типа адреса
	addrType := buf[3]
	log.Debugf("Parse request for %s", string(buf))

	switch addrType {
	case addrTypeIPv4:
		// Если этот тип - IPv4
		address = net.IP(buf[4 : 4+net.IPv4len]).String()
		port = binary.BigEndian.Uint16(buf[4+net.IPv4len : 6+net.IPv4len])
		log.Tracef("parseRequest: IPv4 Port = %d", port)
		break
	case addrTypeDomain:
		// иначе достаем доменное имя
		addrLen := buf[4]
		address = string(buf[5 : 5+addrLen])
		port = binary.BigEndian.Uint16(buf[5+addrLen : 7+addrLen])
		log.Tracef("parseRequest: Domain Port = %d", port)
		break
	default:
		log.Errorf("unsupported address type: %d", addrType)
		return "", 0, errors.New("unsupported address type")
	}
	return address, port, nil
}
