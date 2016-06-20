package evented

import (
	"log"

	"github.com/chuckpreslar/emission"
	"github.com/sunspots/tmi"
)

// Connection wraps a tmi Connection together with an EventEmitter
type Connection struct {
	*emission.Emitter
	tmi.Connector
}

func readLoop(tmi tmi.Connector, connection *Connection) {
	for {
		evt, err := tmi.ReadMessage()
		if err != nil {
			// You could use this to do something useful and
			// handle reconnects, but we're just gonna error out for the sake of simplicity.
			log.Fatal(err)
		} else {
			connection.Emit("message", evt)
			connection.Emit(evt.Command, evt)
		}
	}
}

// New instantiates an EventEmittter
func New(tmi tmi.Connector) *Connection {
	connection := &Connection{emission.NewEmitter(), tmi}
	go readLoop(tmi, connection)
	return connection
}
