package tmi

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Default addresses for WebSocket connections
const (
	wsAddress       = "ws://irc-ws.chat.twitch.tv:80"
	wsSecureAddress = "wss://irc-ws.chat.twitch.tv:443"
)

// WebSocket is a Socket implementation using Web
type WebSocket struct {
	wrMu   sync.Mutex
	config SocketConfig
	socket *websocket.Conn
	secure bool
}

// Open the connection to the Web socket
func (s *WebSocket) Open() (err error) {
	if s.socket != nil {
		s.Close()
	}

	upgradeHeader := http.Header{}
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = s.config.ConnectTimeout

	var addr string
	if s.config.Secure {
		addr = wsSecureAddress
	} else {
		addr = wsAddress
	}

	s.socket, _, err = dialer.Dial(
		addr,
		upgradeHeader,
	)
	s.socket.SetPingHandler(func(message string) error {
		err := s.socket.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(time.Second))
		log.Println("PING " + message)
		if err == websocket.ErrCloseSent {
			return nil
		} else if e, ok := err.(net.Error); ok && e.Temporary() {
			return nil
		}
		return err
	})
	return err
}

// Close the connection to the Web socket
func (s *WebSocket) Close() {
	s.socket.Close()
	s.socket = nil
}

// ReadLine reads a line from the socket, with read timeout
func (s *WebSocket) ReadLine() ([]byte, error) {
	s.socket.SetReadDeadline(time.Now().Add(s.config.ReadTimeout))
	t, b, err := s.socket.ReadMessage()
	s.socket.SetReadDeadline(time.Time{})
	if err != nil {
		return b, err
	}
	if t != websocket.TextMessage {
		return b, fmt.Errorf("Read non-text message of type: %d: %v", t, b)
	}
	return b, err
}

// WriteLine writes a line to the socket, with read timeout
func (s *WebSocket) WriteLine(m []byte) error {
	s.wrMu.Lock()
	defer s.wrMu.Unlock()
	s.socket.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))
	err := s.socket.WriteMessage(websocket.TextMessage, m)
	s.socket.SetWriteDeadline(time.Time{})
	return err
}

// NewWebSocket  returns a Socket using WebSocket with default settings
func NewWebSocket(secure bool) *WebSocket {
	return &WebSocket{
		config: SocketConfig{
			Secure:         secure,
			ConnectTimeout: defaultConnectTimeout,
			ReadTimeout:    defaultReadTimeout,
			WriteTimeout:   defaultWriteTimeout,
		},
	}
}
