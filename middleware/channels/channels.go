// Package channels is a plugin for handling multiple channels on a single connection
// managing joining, leaving, currently joined channels, etc.
package channels

import (
	"fmt"

	"github.com/sunspots/tmi"
)

// Channel contains options and state for an individual channel
type Channel struct {
	group     *Group
	name      string
	in        bool              // Track wether we have successfully joined.
	UserState map[string]string // Connected user's userstate
}

// Group object for containing and managing channels
type Group struct {
	Channels map[string]*Channel
	Conn     *tmi.Connection
}

func (ch *Channel) chanHandler(m *tmi.Message) {
	//USERSTATE data is saved in the channel for later.
	switch m.Command {
	case "USERSTATE":
		if ch.group.Conn.Username == ch.name[1:] {
			m.Tags["user_type"] = "broadcaster"
		}
		ch.UserState = m.Tags
	case "366":
		ch.in = true
		fmt.Printf("Joined channel %s successfully\n", m.Params[0])
	}
}

// MiddleWare is a pluggable intermediate on message handling for the TMI construct.
func (chs *Group) MiddleWare(m *tmi.Message, err error) (*tmi.Message, error) {
	if err != nil {
		return m, err
	}

	if ch, ok := chs.Channels[m.Params[0]]; ok {
		ch.chanHandler(m)
	} else {
		//Do something if handling a channel not joined
	}
	return m, nil
}

// Reset the Group object
func (chs *Group) Reset() {
	chs.Channels = make(map[string]*Channel)
}

// Join a specified channel
func (chs *Group) Join(c string) *Channel {
	if c[0] != '#' {
		c = "#" + c
	}
	if _, ok := chs.Channels[c]; !ok {
		chs.Channels[c] = &Channel{
			group: chs,
			name:  c,
		}
		chs.Conn.Join(c)
	} else {
		fmt.Printf("%s already joined \n", c)
	}
	return chs.Channels[c]
}

// Part a specified channel, if the channel isn't joined, it fails silently
func (chs *Group) Part(c string) {
	if c[0] != '#' {
		c = "#" + c
	}
	if _, ok := chs.Channels[c]; ok {
		chs.Conn.Send("PART " + c)
		delete(chs.Channels, c)
	}
}

// In returns a boolean describing wether the channel is currently joined or not
func (chs *Group) In(c string) bool {
	if c[0] != '#' {
		c = "#" + c
	}
	_, ok := chs.Channels[c]
	return ok
}

// New creates a new channel group to manage channels on top of the TMI connection
func New(conn *tmi.Connection) *Group {
	chs := &Group{
		Channels: make(map[string]*Channel),
		Conn:     conn,
	}
	return chs
}
