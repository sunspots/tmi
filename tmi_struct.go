package tmi

import (
	"net"
	"sync"
	"time"
)

//Connection is the main struct for for containing an active connection
type Connection struct {
	sync.WaitGroup
	sync.Mutex
	Server      string // Server to connect to
	Port        string // Port to connect to
	Nick        string // Nick to connect with
	Token       string // OAuth token, without "oauth:" prefix
	TC          string // TWITCHCLIENT string to use on connection
	stopped     bool
	end         chan bool
	send        chan string
	Debug       bool // Debug decides if debug messages should be printed.
	Error       chan error
	socket      net.Conn
	MessageChan chan *Message
	Timeout     time.Duration
	KeepAlive   time.Duration
	lastMessage time.Time
}
