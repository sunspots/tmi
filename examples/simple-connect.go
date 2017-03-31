package main

import (
	"log"

	"github.com/sunspots/tmi"
)

const (
	channel = "#sunspots"
)

func readLoop(connection *tmi.Connection) {
	for {
		evt, err := connection.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if evt.Command == "PRIVMSG" {
			log.Printf("%s says %s to %s: %s\n", evt.From, evt.Command, evt.Channel(), evt.Trailing)
		}

	}
}

func main() {
	connection := tmi.Anonymous() // Initialise the connection object with an anonymous login, as well as logging in anonymously
	connection.Debug = true       // Prints out raw incoming and outgoing messages
	connection.Join(channel)
	readLoop(connection)
	// readLoop returns in case of disconnect
}
