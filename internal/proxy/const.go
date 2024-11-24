package proxy

// Constants for SOCKS5
const (
	socks5Version = 0x05
	cmdConnect    = 0x01

	addrTypeIPv4   = 0x01
	addrTypeDomain = 0x03

	replySuccess = 0x00
	replyFailure = 0x01
)
