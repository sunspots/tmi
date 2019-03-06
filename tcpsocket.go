package tmi

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"
)

// Default addresses for TCP connections
const (
	tcpAddress       = "irc.chat.twitch.tv:6667"
	tcpSecureAddress = "irc.chat.twitch.tv:6697"
)

// TCPSocket is a Socket implementation over TCP, with optional TLS.
type TCPSocket struct {
	wrMu   sync.Mutex // Serialize write operations
	config SocketConfig
	conn   net.Conn
	reader *bufio.Reader
}

// Open the connection to the TCP socket
func (s *TCPSocket) Open() (err error) {
	if s.conn != nil {
		s.Close()
	}
	dialer := net.Dialer{Timeout: s.config.ConnectTimeout}
	if s.config.Secure {
		s.conn, err = tls.DialWithDialer(&dialer, "tcp", tcpSecureAddress, &tls.Config{})
	} else {
		s.conn, err = dialer.Dial("tcp", tcpAddress)
	}
	if err != nil {
		return err
	}

	s.reader = bufio.NewReaderSize(s.conn, maxMessageSize)
	return err
}

// Close the connection to the TCP socket
func (s *TCPSocket) Close() {
	s.conn.Close()
	s.conn = nil
}

// ReadLine reads a line from the socket, with read timeout
func (s *TCPSocket) ReadLine() ([]byte, error) {
	s.conn.SetReadDeadline(time.Now().Add(s.config.ReadTimeout))
	b, err := s.reader.ReadBytes('\n')
	s.conn.SetReadDeadline(time.Time{})
	return b, err
}

// WriteLine writes a line to the socket, with write timeout
func (s *TCPSocket) WriteLine(m []byte) error {
	s.wrMu.Lock()
	defer s.wrMu.Unlock()
	s.conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))
	_, err := fmt.Fprintf(s.conn, "%s\r\n", m)
	s.conn.SetWriteDeadline(time.Time{})
	return err
}

// NewTCPSocket  returns a Socket using TCP with default settings
func NewTCPSocket(secure bool) *TCPSocket {
	return &TCPSocket{
		config: SocketConfig{
			Secure:         secure,
			ConnectTimeout: defaultConnectTimeout,
			ReadTimeout:    defaultReadTimeout,
			WriteTimeout:   defaultWriteTimeout,
		},
	}
}
