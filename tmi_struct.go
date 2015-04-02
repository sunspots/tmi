package tmi

import (
	"net"
	"sync"
	"time"
)

//Connection is the main struct for for containing an active connection
type Connection struct {
	sync.WaitGroup
	sync.Mutex
	Server      string // Server to connect to
	Port        string // Port to connect to
	Nick        string // Nick to connect with
	Token       string // OAuth token, without "oauth:" prefix
	TC          string // TWITCHCLIENT string to use on connection
	stopped     bool
	end         chan bool
	send        chan string
	Debug       bool // Debug decides if debug messages should be printed.
	Error       chan error
	socket      net.Conn
	EventChan   chan *Event
	Timeout     time.Duration
	KeepAlive   time.Duration
	lastMessage time.Time
}

// Event is a struct that represents events recieved from the server
type Event struct {
	From    string            `json:"from"`
	Command string            `json:"command"`
	Params  []string          `json:"params"`
	RawTags map[string]string `json:"rawTags"`
	Tags    *UserTags         `json:"tags"`
	Time    time.Time         `json:"time"`
	Cleared bool              `json:"cleared, omitempty"`
	Action  bool              `json:"action, omitempty"`
}

// Message returns the computed message of the event, or an empty string
func (evt *Event) Message() string {
	if len(evt.Params) == 0 {
		return ""
	}
	return evt.Params[len(evt.Params)-1]
}

// Channel returns the computed target channel of the event, if there is one
func (evt *Event) Channel() string {
	if len(evt.Params) < 1 {
		return ""
	}
	if evt.Params[0][0] == '#' {
		return evt.Params[0]
	}
	return ""
}

// UserTags object keeps an object that contains parsed IRCv3 tags
type UserTags struct {
	Color      string   `json:"color"`
	Emotes     []*Emote `json:"emotes"`
	EmoteSets  []int    `json:"emotesets"`
	Subscriber bool     `json:"subscriber"`
	Turbo      bool     `json:"turbo"`
	UserType   string   `json:"user_type"`
}

// Emote struct for storing an emote, with a single from/to position
type Emote struct {
	Id   string `json:"id"`
	From int    `json:"from"`
	To   int    `json:"to"`
}

//Sorting interface for Emotes
type ByPos []*Emote

func (a ByPos) Len() int           { return len(a) }
func (a ByPos) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPos) Less(i, j int) bool { return a[i].From < a[j].From }
