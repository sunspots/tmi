package tmi_test

import (
	"reflect"
	"testing"

	"github.com/sunspots/tmi"
)

var (
	// Used to check messages received from the "server"
	serverMessages = []tmi.Message{
		tmi.Message{
			Raw:     ":sunspots PRIVMSG #sunspots :test message, lol\n",
			From:    "sunspots",
			Command: "PRIVMSG",
			Params:  []string{"#sunspots"},
			Channel: "sunspots",
			Text:    "test message, lol",
		},
		tmi.Message{
			Raw:     ":sunspots PRIVMSG #sunspots :test message, lol\n",
			From:    "sunspots",
			Command: "PRIVMSG",
			Params:  []string{"#sunspots"},
			Channel: "sunspots",
			Text:    "test message, lol",
		},
		tmi.Message{
			Raw:     ":sunspots PRIVMSG #sunspots :test message, lol\n",
			From:    "sunspots",
			Command: "PRIVMSG",
			Params:  []string{"#sunspots"},
			Channel: "sunspots",
			Text:    "test message, lol",
		},
		tmi.Message{
			Raw:     ":sunspots PRIVMSG #sunspots :test message, lol\n",
			From:    "sunspots",
			Command: "PRIVMSG",
			Params:  []string{"#sunspots"},
			Channel: "sunspots",
			Text:    "test message, lol",
		},
	}
	// TODO: Check that the client sends the expected messages to the "server"
	clientMessages = []string{
		"hello",
		"world",
	}
)

// Test receiving messages from the server
func TestRecieve(t *testing.T) {
	socket := tmi.NewPipeSocket()
	client := tmi.New("foo", "bar", socket)
	client.Connect()

	go func() {
		for _, msg := range serverMessages {
			socket.InputMessage(msg)
		}
		socket.Close()
	}()

	for _, exp := range serverMessages {
		read, err := client.ReadMessage()
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(exp, read) {
			t.Errorf("Expected: '%v', Got: '%v'", exp, read)
		}
	}
}

// Test sending messages to the server
func TestSend(t *testing.T) {
	socket := tmi.NewPipeSocket()
	client := tmi.New("foo", "bar", socket)
	client.Connect()

	go func() {
		for _, msg := range clientMessages {
			client.Send(msg)
		}
	}()

	for _, exp := range clientMessages {
		read, err := socket.Output()
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(exp, read) {
			t.Errorf("Expected: '%v', Got: '%v'", exp, read)
		}
	}
}
