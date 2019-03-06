package tmi

import (
	"fmt"
	"strings"
)

// The Client encapsulates state and methods for communicating over a socket.
// Most use cases are covered by calling Anonymous, Default or New
type Client struct {
	Username string
	Token    string
	socket   Socket
}

// Connect to the TMI server with the socket
func (c *Client) Connect() error {
	err := c.socket.Open()
	if err != nil {
		return err
	}
	if len(c.Token) != 0 {
		c.Send("PASS oauth:" + c.Token)
	}

	c.Send("NICK " + c.Username)
	c.Send("CAP REQ :twitch.tv/tags twitch.tv/commands")
	return nil
}

// ReadMessage reads from the TMI server, blocking call.
func (c *Client) ReadMessage() (msg Message, err error) {
	line, err := c.socket.ReadLine()
	if err != nil {
		return msg, err
	}
	err = parseMessage(string(line), &msg)
	if err != nil {
		return msg, err
	}
	if msg.Command == "PING" {
		c.SendMessage(Message{Command: "PONG"})
		return c.ReadMessage()
	}
	return msg, err
}

// Say sends a regular chat message to a target channel (exclude the # from channel name)
func (c *Client) Say(channel string, text string) error {
	return c.Sendf("PRIVMSG #%s :%s", strings.ToLower(channel), text)
}

// Join a channel (exclude the # from channel name)
func (c *Client) Join(channel string) error {
	return c.Sendf("Join #%s", channel)
}

// SendBytes sends a bytearray to the TMI server
func (c *Client) SendBytes(m []byte) error {
	return c.socket.WriteLine(m)
}

// Send convenience for sending strings (see Send)
func (c *Client) Send(m string) error {
	return c.SendBytes([]byte(m))
}

// Sendf acts like fmt.Sprintf to send a formatted string with arguments
func (c *Client) Sendf(m string, args ...interface{}) error {
	return c.SendBytes([]byte(fmt.Sprintf(m, args...)))
}

// SendMessage conveniently encodes your message and sends it
func (c *Client) SendMessage(message Message) error {
	return c.SendBytes(message.Bytes())
}

// New returns a new client with your specified socket
func New(username string, token string, socket Socket) *Client {
	return &Client{
		socket:   socket,
		Username: username,
		Token:    token,
	}
}

// Default returns a new client with default socket settings (TCP/Secure)
func Default(username string, token string) *Client {
	socket := NewTCPSocket(true)
	return New(username, token, socket)
}

// Anonymous client connects without credentials, sent messages are dropped.
func Anonymous() *Client {
	return Default(AnonymousUser, "")
}
