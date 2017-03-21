// Package tmi is a simple library for connecting to Twitch's pseudo-IRC servers, aka Twitch Message Interface
package tmi

import (
	"bufio"
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

// ReadMessage reads an incoming message from the server, blocking until a message is recieved
func (tmi *Connection) ReadMessage() (*Message, error) {
	evt, ok := <-tmi.MessageChan
	var err error
	if !ok {
		err = errors.New("read message channel closed")
	} else {
		err = nil
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
		//tmi.Error <- errors.New("Disconnect Called")
	}
}

// Reconnect to a connected server
func (tmi *Connection) Reconnect() error {
	tmi.Disconnect()
	return tmi.Connect()
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

// The main loop to read messages from the server, runs as a goroutine.
func (tmi *Connection) readLoop() {
	defer func() {
		log.Println("Reader closed")
		tmi.Done()
	}()
	br := bufio.NewReaderSize(tmi.socket, maxMessageSize)

	for {
		select {
		case <-tmi.end:
			return
		default:
			if tmi.socket != nil {
				tmi.socket.SetReadDeadline(time.Now().Add(tmi.Timeout*2 + tmi.KeepAlive))
			}
			msg, err := br.ReadString('\n')
			if err != nil {
				tmi.Error <- err
				return
			}
			if tmi.socket != nil {
				var zero time.Time
				tmi.socket.SetReadDeadline(zero)
			}
			tmi.lastMessage = time.Now()

			if tmi.Debug {
				dbg.Print("< ", msg)
			}

			message := ParseMessage(msg)
			if message.Command == "PING" {
				tmi.Send("PONG " + message.Trailing)
				continue
			}
			tmi.MessageChan <- message
		}
	}
}

func (tmi *Connection) writeLoop() {
	defer func() {
		log.Println("Writer closed")
		tmi.Done()
	}()
	for {
		select {
		case s, ok := <-tmi.send:
			if !ok {
				tmi.Error <- errors.New("send channel closed")
				return
			}
			if tmi.socket == nil {
				tmi.Error <- errors.New("no socket to write to")
				return
			}
			if s == "" {
				continue
			}
			if tmi.Debug {
				dbg.Printf("> %s\n", s)
			}
			tmi.socket.SetWriteDeadline(time.Now().Add(tmi.Timeout))
			_, err := fmt.Fprintf(tmi.socket, "%s\r\n", s)
			var zero time.Time
			tmi.socket.SetWriteDeadline(zero)

			if err != nil {
				tmi.Error <- err
				return
			}
		case <-tmi.end:
			return
		}
	}
}

func (tmi *Connection) pingLoop() {
	defer func() {
		log.Println("Pinger stopped")
		tmi.Done()
	}()
	ticker := time.NewTicker(tmi.Timeout) // Tick for monitoring
	for {
		select {
		case _, ok := <-ticker.C:
			if !ok {
				// The ticker has been closed, probably shouldn't happen before tmi.end is closed
				log.Println("ticker not ok")
				return
			}

			//Ping if we haven't received anything from the server within the keep alive period
			if time.Since(tmi.lastMessage) >= tmi.KeepAlive {
				tmi.Sendf("PING %d", time.Now().UnixNano())
			}
		case <-tmi.end:
			ticker.Stop()
			return
		}
	}
}

func (tmi *Connection) controlLoop() {
	for !tmi.Stopped() {
		err := <-tmi.Error
		if tmi.Stopped() {
			// If one of the goroutines errors when we've already stopped,
			// There's no need to handle the error.
			continue
		}
		log.Printf("Error, disconnecting: %s\n", err)
		tmi.Done()
		log.Println("controlLoop done, waiting for disconnect")
		tmi.Disconnect()
	}
	log.Println("controlloop terminated")
}

// Join joins a channel
func (tmi *Connection) Join(channel string) {
	tmi.Send("JOIN " + channel)
}

// Connect connects to the server, starts reading and authenticates
func (tmi *Connection) Connect() (err error) {

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

	tmi.Send("NICK " + tmi.Nick)
	tmi.Send("CAP REQ :twitch.tv/tags twitch.tv/commands")

	return nil
}

// New returns a new connection object, ready to connect
func New(nick, token string) *Connection {
	tmi := &Connection{
		Server:    "irc.chat.twitch.tv",
		Port:      "6667",
		Debug:     false,
		socket:    nil,
		stopped:   true,
		Timeout:   15 * time.Second,
		KeepAlive: 30 * time.Second,
		Error:     make(chan error, 3),
		Nick:      nick,
		Token:     token,
	}
	return tmi
}

// Connect is a shortcut to creating a new connection, connecting to it and returns the connection
func Connect(nick, token string) *Connection {
	new := New(nick, token)
	new.Connect()
	return new
}

func Anonymous() *Connection {
	new := New("justinfan999", "")
	new.Connect()
	return new
}
