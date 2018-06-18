// Package tmi is a simple library for connecting to Twitch's pseudo-IRC servers, aka Twitch Message Interface
package tmi

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	dbg = log.New(os.Stdout, "TMI: ", log.LstdFlags)
)

const (
	// Maximum message size allowed from peer.
	maxMessageSize = 512
	// Default timeout
	timeout = 15 * time.Second
	// Default keepalive
	keepAlive = 30 * time.Millisecond
)

// Send sends messages to the TMI server
func (tmi *Connection) Send(s string) {
	if !tmi.Stopped() {
		tmi.send <- s
	} else {
		dbg.Printf("unable to send %s on closed connection \n", s)
	}
}

// Sendf sends a message, with format and params, wrapper around fmt.Sprintf
func (tmi *Connection) Sendf(format string, a ...interface{}) {
	tmi.Send(fmt.Sprintf(format, a...))
}

// Join simply sends a join message
func (tmi *Connection) Join(channel string) {
	tmi.Send("JOIN " + channel)
}

// Stopped tells us wether the client is stopped or not
func (tmi *Connection) Stopped() bool {
	tmi.Lock()
	defer tmi.Unlock()
	return tmi.stopped
}

// Update the stopped status when starting or stopping the client
func (tmi *Connection) setStopped(value bool) {
	tmi.Lock()
	tmi.stopped = value
	tmi.Unlock()
}

// ReadMessage reads an incoming message from the server,
// blocking until a message is recieved or an error occurs
func (tmi *Connection) ReadMessage() (*Message, error) {
	evt, ok := <-tmi.MessageChan
	var err error
	if !ok {
		err = errors.New("read message channel closed")
	}
	return evt, err
}

// Disconnect from the server
func (tmi *Connection) Disconnect() {
	if !tmi.Stopped() {
		tmi.setStopped(true)
		close(tmi.end)
		tmi.Wait()
		close(tmi.send)
		tmi.socket.Close()
		tmi.socket = nil
		close(tmi.MessageChan)
		dbg.Println("Disconnected!")
	}
}

// Reconnect to a connected server
func (tmi *Connection) Reconnect() error {
	tmi.Disconnect()
	return tmi.Connect()
}

// Connect connects to the server, starts routines and authenticates
func (tmi *Connection) Connect() (err error) {
	if !tmi.stopped {
		return errors.New("Can't attempt to Connect with a Connection that isn't stopped!")
	}

	if len(tmi.Token) > 0 && tmi.Token[0:6] == "oauth:" {
		tmi.Token = tmi.Token[6:]
	}

	tmi.socket, err = net.DialTimeout("tcp", tmi.Server+":"+tmi.Port, tmi.Timeout)
	if err != nil {
		return err
	}

	log.Printf("Connected to TMI server %s (%s)\n", tmi.Server, tmi.socket.RemoteAddr())
	tmi.Lock()
	tmi.end = make(chan bool)
	tmi.send = make(chan string, 10)
	tmi.MessageChan = make(chan *Message, 50)
	tmi.stopped = false
	tmi.Unlock()
	tmi.Add(4)
	go tmi.readLoop()
	go tmi.writeLoop()
	go tmi.pingLoop()
	go tmi.controlLoop()

	//Authenticate

	if len(tmi.Token) != 0 {
		tmi.Send("PASS oauth:" + tmi.Token)
	}

	tmi.Send("NICK " + tmi.Username)
	tmi.Send("CAP REQ :twitch.tv/tags twitch.tv/commands")

	return nil
}

// New returns a new connection object, ready to connect
func New(username, token string) *Connection {
	tmi := &Connection{
		Server:    "irc.chat.twitch.tv",
		Port:      "6667",
		Debug:     false,
		socket:    nil,
		stopped:   true,
		Timeout:   timeout,
		KeepAlive: keepAlive,
		Error:     make(chan error, 3),
		Username:  username,
		Token:     token,
	}
	return tmi
}

// Connect is a shortcut to calling New and connecting before returning
func Connect(username, token string) *Connection {
	new := New(username, token)
	new.Connect()
	return new
}

// Anonymous is a shortcut to calling New and connecting, without personal credentials, before returning
func Anonymous() *Connection {
	new := New("justinfan999", "")
	new.Connect()
	return new
}
