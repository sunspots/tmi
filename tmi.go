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
func (tmi *Connection) Sendf(format string, a ...interface{}) {
	tmi.Send(fmt.Sprintf(format, a...))
}
func (tmi *Connection) ReadEvent() (*Event, error) {
	evt, ok := <-tmi.EventChan
	var err error
	if !ok {
		err = errors.New("read event channel closed")
	} else {
		err = nil
	}
	return evt, err
}

func (tmi *Connection) Disconnect() {
	if !tmi.Stopped() {
		tmi.SetStopped(true)
		close(tmi.end)
		tmi.Wait()
		close(tmi.send)
		tmi.socket.Close()
		tmi.socket = nil
		close(tmi.EventChan)
		dbg.Println("Disconnected!")
		//tmi.Error <- errors.New("Disconnect Called")
	}
}

func (tmi *Connection) Reconnect() error {
	tmi.Disconnect()
	return tmi.Connect(tmi.Nick, tmi.Token)
}

func (tmi *Connection) Stopped() bool {
	tmi.Lock()
	defer tmi.Unlock()
	return tmi.stopped
}

func (tmi *Connection) SetStopped(value bool) {
	tmi.Lock()
	tmi.stopped = value
	tmi.Unlock()
}

// The main loop to read messages from the server, runs as a goroutine.
func (tmi *Connection) readLoop() {
	defer func() {
		log.Println("readloop done")
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

			event := parseEvent(msg)
			if event.Command == "PING" {
				tmi.Send("PONG " + event.Message())
				continue
			}
			tmi.EventChan <- event
		}
	}
}

func (tmi *Connection) writeLoop() {
	defer func() {
		log.Println("writeloop done")
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
		log.Println("pingloop done")
		tmi.Done()
	}()
	ticker := time.NewTicker(tmi.Timeout) // Tick for monitoring
	for {
		select {
		case _, ok := <-ticker.C:
			if !ok {
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
			log.Println("TMI Error thrown while stopped, should not happen!")
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
func (tmi *Connection) Connect(nick, token string) (err error) {
	if len(nick) == 0 {
		return errors.New("empty 'nick'")
	}
	if len(token) == 0 {
		return errors.New("empty 'token'")
	}
	tmi.Nick = nick
	tmi.Token = token

	if tmi.Token[0:6] == "oauth:" {
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
	tmi.EventChan = make(chan *Event, 50)
	tmi.Unlock()
	tmi.SetStopped(false)
	tmi.Add(4)
	go tmi.readLoop()
	go tmi.writeLoop()
	go tmi.pingLoop()
	go tmi.controlLoop()

	//Authenticate
	tmi.Send("PASS oauth:" + tmi.Token)
	tmi.Send("NICK " + tmi.Nick)
	tmi.Send("CAP REQ :twitch.tv/tags twitch.tv/commands")
	if len(tmi.TC) > 1 {
		tmi.Send(tmi.TC)
	}
	//Answer pings

	return nil
}

// New returns a new connection object, ready to connect
func New() *Connection {
	tmi := &Connection{
		Server:    "irc.twitch.tv",
		Port:      "6667",
		Debug:     false,
		socket:    nil,
		stopped:   true,
		Timeout:   15 * time.Second,
		KeepAlive: 30 * time.Second,
		Error:     make(chan error, 3),
	}
	return tmi
}
