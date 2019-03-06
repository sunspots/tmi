package tmi

import "time"

const (
	// Default timeout for sockets' opening
	defaultConnectTimeout = 2 * time.Minute
	// Default timeout for reading a message from the server
	defaultReadTimeout = 6 * time.Minute
	// Default timeout when trying to write a message to the server
	defaultWriteTimeout = 2 * time.Minute
)

// Socket wraps the socket implementation with standardised methods
type Socket interface {
	// Opens the connection to the server
	Open() error
	// Closes the connection to the server
	Close()
	// Reads a line, unlikely to be thread-safe
	ReadLine() ([]byte, error)
	// Write line to the server
	// should be thread-safe/synchronized to allow async use
	WriteLine([]byte) error
}

// SocketConfig for commmon settings used by the default socket types
type SocketConfig struct {
	Secure         bool          // Request the socket to use TLS or other secure transport
	ConnectTimeout time.Duration // Timeout for opening the socket
	ReadTimeout    time.Duration // Timeout for reading from the socket
	WriteTimeout   time.Duration // Timeout when writing to the socket
}
