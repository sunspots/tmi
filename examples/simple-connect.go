package main

import (
	"bufio"
	"fmt"
	"github.com/sunspotseu/tmi"
	"log"
	"os"
)

var exampleChannel = "#sunspots"

func readLoop(connection *tmi.Connection) {
	for {
		evt, err := connection.ReadEvent()
		if err != nil {
			// You could use this to do something useful and
			// handle reconnects, but we're just gonna error for the sake of simplicity.
			log.Fatal(err)
		} else {
			if evt.Command == "PRIVMSG" {
				log.Printf("%s says %s to %s: %s\n", evt.From, evt.Command, evt.Channel(), evt.Message())
			} else {
				log.Println("Unhandled event", evt)
			}
		}
	}
}

func main() {
	connection := tmi.New() // Initialise the connection object with required login
	connection.Debug = true // Prints out raw incoming and outgoing messages

	connection.Connect("sunspots", "oauth:9m48oiumz0bz3v8ec4g6ph8ylw4roe") // Connect and authenticate

	connection.Join(exampleChannel)
	go readLoop(connection)

	// Now, since it's running things in a goroutine, we don't want main to exit yet,
	// you can do what you want here, but maybe something like...
	// Take input to send messages
	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')
		connection.Send(fmt.Sprintf("PRIVMSG %s %s", exampleChannel, input))
	}

}
