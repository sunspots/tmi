package main

import (
	"log"

	"github.com/sunspots/tmi"
)

const (
	channel = "lirik"
)

func main() {

	socket := tmi.NewWebSocket(true)
	client := tmi.New(tmi.AnonymousUser, "", socket)

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	client.Join(channel)

	for {
		m, err := client.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(m.Raw)
	}
}
