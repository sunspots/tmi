package tmi

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"time"
)

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

// The writeloop synchronously sends messages
// that it picks up from the send channel
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

// The pingloop sends automatic, slightly higher frequency, pings
// to allow a shorter read timeout and disconnect detection.
// It skips unnecessary pings in case the last received message
// is within the KeepAlive timeframe
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

// The control loop manages disconnection if one of the other loops
// sends an error on the tmi.Error channel
// That way, the loops can just send an error and quit without caring about other routines
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
