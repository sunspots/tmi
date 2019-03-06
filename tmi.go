// Package tmi implements a TMI client (IRC with Twitch flavour) to connect to the Twitch chat.
//
// Overview
//
// The Client type represents the client connection to the Twitch chat.
// Client objects can be created manually, but the most common use cases are covered by calling Anonymous, Default or New.
// A client calls client.Connect() in order to establish the connection.
//
// 	client := tmi.Default("username", "12345")
// 	err := client.Connect()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	client.Join("twitchdev")
// 	.. use client to send and receive messages
//
// Unless you use Anonymous or Default, you need to supply a socket implementation,
// the default socket implementations are TCP and WebSocket, both come with optional (but recommended) TLS.
//
// 	socket := tmi.NewWebSocket(true)
// 	client := tmi.New("username", "12345", socket)
//
// For most use cases, call the client's Say and ReadMessage methods to communicate with chat.
// For sending, there are also the Send(string), Sendf(like fmt.Printf) SendBytes(Bytearray) and SendMessage(*Message) methods.
// This example shows how to automatically respond to all messages:
//
//	for {
//		m, err := client.ReadMessage()
//		if err != nil {
//			log.Fatal(err)
//		}
//		client.Say(m.Channel, "Hello " + m.From)
//	}
//
// Detecting dropped connections
//
// When working with real-time chat, missing out on messages sucks.
// In case a connection drops, we want to detect it early and attempt to reconnect.
// The main method for detecting a dropped connection is to attempt communication.
// Twitch checks that your connection is alive by sending a PING message every ~5 minutes, which `tmi` responds to automatically.
// The application must read on the connection in order to process PINGs (ReadMessage will block until it's handled, and attempt to return the next message)
//
// If you want to avoid missing messages, it could be worth sending your own PING messages to Twitch.
// Usually, sending the message regularly is good enough (no need to check for PONG).or other control messages (RECONNECT, etc.)
// A snippet like this might be enough.
//
// 	stop := new(chan struct{})
// 	ticker := time.NewTicker(time.Minute)
// 	go func() {
// 		for {
// 			select {
// 			case <-ticker.C:
// 				err := c.Sendf("PING %d", time.Now().UnixNano())
// 				if err != nil {
// 					return
// 				}
// 			case <-stop:
// 				ticker.Stop()
// 				return
// 			}
// 		}
// 	}()
// close(stop)
package tmi

const (
	// Maximum message size allowed from peer.
	maxMessageSize = 512
	// AnonymousUser is just a username you can use to connect anonymously
	AnonymousUser = "justinfan999"
)
