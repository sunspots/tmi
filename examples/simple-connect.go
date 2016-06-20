package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/sunspots/tmi"
)

const (
	User    = "sunspots"
	OAuth   = "oauth:9m48oiumz0bz3v8ec4g6ph8ylw4roe"
	Channel = "#sunspots"
)

func readLoop(connection *tmi.Connection) {
	for {
		evt, err := connection.ReadMessage()
		if err != nil {
			// You could use this to do something useful and
			// handle reconnects, but we're just gonna error out for the sake of simplicity.
			log.Fatal(err)
		} else {
			if evt.Command == "PRIVMSG" {
				log.Printf("%s says %s to %s: %s\n", evt.From, evt.Command, evt.Channel(), evt.Trailing)
			} else {
				log.Println("Unhandled event", evt)
			}
		}
	}
}

func main() {
	connection := tmi.New() // Initialise the connection object with required login
	connection.Debug = true // Prints out raw incoming and outgoing messages

	connection.Connect(User, OAuth) // Connect and authenticate

	connection.Join(Channel)
	go readLoop(connection)

	// Now, since it's running things in a goroutine, we don't want main to exit yet,
	// you can do what you want here, but maybe something like...
	// Take input to send messages
	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')
		connection.Send(fmt.Sprintf("PRIVMSG %s %s", Channel, input))
	}

}
