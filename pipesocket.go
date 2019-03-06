package tmi

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// PipeSocket is a Socket implementation using an in-memory pipe
type PipeSocket struct {
	wrMu         sync.Mutex // Serialize writes (pipe socket is synchronised, so this is redundant)
	config       SocketConfig
	server       net.Conn
	client       net.Conn
	clientReader *bufio.Reader
	serverReader *bufio.Reader
}

// Open the connection to the TCP socket
func (s *PipeSocket) Open() (err error) {
	s.server, s.client = net.Pipe()

	s.clientReader = bufio.NewReaderSize(s.client, maxMessageSize)
	s.serverReader = bufio.NewReaderSize(s.server, maxMessageSize)
	return err
}

// Close the connection to the TCP socket
func (s *PipeSocket) Close() {
	s.server.Close()
	s.client.Close()
}

// ReadLine reads a single line from the pipe
func (s *PipeSocket) ReadLine() ([]byte, error) {
	s.client.SetReadDeadline(time.Now().Add(s.config.ReadTimeout))
	b, err := s.clientReader.ReadBytes('\n')
	s.client.SetReadDeadline(time.Time{})
	return b, err
}

// WriteLine writes a line to the pipe
func (s *PipeSocket) WriteLine(m []byte) error {
	s.wrMu.Lock()
	defer s.wrMu.Unlock()
	s.client.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))
	_, err := fmt.Fprintf(s.client, "%s\n", m)
	s.client.SetWriteDeadline(time.Time{})
	return err
}

// Input some new data to the input pipe! Emulates the server.
func (s *PipeSocket) Input(p []byte) {
	s.server.Write(p)
}

// InputMessage formats and uses Input (see Input)
func (s *PipeSocket) InputMessage(m Message) {
	s.Input(m.Bytes())
	s.Input([]byte("\n"))
}

// Output gets the output from the server, aka the sent messages
func (s *PipeSocket) Output() (string, error) {
	b, err := s.serverReader.ReadBytes('\n')
	return strings.TrimSpace(string(b)), err
}

// NewPipeSocket returns a Socket using TCP with default settings
func NewPipeSocket() *PipeSocket {
	return &PipeSocket{
		config: SocketConfig{
			ReadTimeout:  100 * time.Millisecond,
			WriteTimeout: 100 * time.Millisecond,
		},
	}
}
