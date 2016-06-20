package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/sunspots/tmi"
	"github.com/sunspots/tmi/middleware/evented"
)

const (
	user    = "sunspots"
	oAuth   = "oauth:foobar"
	channel = "#lirik"
)


func main() {
	connection := evented.New(tmi.Connect(user, oAuth)) // Initialise the connection object
	connection.Join(channel)

  connection.On("PRIVMSG", func(msg *tmi.Message) {
    log.Println(msg.From + ": " + msg.Trailing)
  })
	// Now, since it's running things in a goroutine, we don't want main to exit yet,
	// you can do what you want here, but maybe something like...
	// Take input to send messages
	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')
		connection.Send(fmt.Sprintf("PRIVMSG %s %s", channel, input))
	}
}
