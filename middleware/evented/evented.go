package evented

import (
	"log"

	"github.com/chuckpreslar/emission"
	"github.com/sunspots/tmi"
)
// EventedConnection wraps a tmi Connection together with an EventEmitter
type EventedConnection struct {
  *emission.Emitter
  *tmi.Connection
}

func readLoop(connection *tmi.Connection, emitter *EventedConnection) {
	for {
		evt, err := connection.ReadMessage()
		if err != nil {
			// You could use this to do something useful and
			// handle reconnects, but we're just gonna error out for the sake of simplicity.
			log.Fatal(err)
		} else {
      emitter.Emit("message", evt)
      emitter.Emit(evt.Command, evt)
		}
	}
}

// New instantiates an EventEmittter
func New(tmi *tmi.Connection) *EventedConnection {
	emitter := &EventedConnection{emission.NewEmitter(), tmi}
  go readLoop(tmi, emitter)
	return emitter
}
