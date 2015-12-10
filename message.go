package tmi

import (
	"bytes"
	"sort"
	"strconv"
	"strings"
)

const (
	actionPrefixLen = 8             // Length of the action start
	actionPrefix    = "\x01ACTION " // Indicates the start of an ACTION(/me) message
	actionSuffix    = '\x01'        // Indicates the end of the same ^

	prefix     byte = 0x3A // ":" Prefix/trailing/emoticon-separator
	prefixUser byte = 0x21 // "!" Username
	prefixTags byte = 0x40 // "@" Tags (and hostname, but we don't care about hostnames in tmi)
	space      byte = 0x20 // " " Separator

	tagSep byte = ';' // Tags separator
	tagAss byte = '=' // Tags assignator
	emSep  byte = '-' // Emote separator

	maxLength = 510 // Maximum length is 512 - 2 for the line endings.
)

// Message struct contains all the relevant data for a message
type Message struct {
	From     string            `json:"from"`
	Command  string            `json:"command"`
	Params   []string          `json:"params"`
	Trailing string            `json:"trailing"`
	Tags     map[string]string `json:"tags"`
	Emotes   []*Emote          `json:"emotes"`
	Action   bool              `json:"action, omitempty"`
}

// Emote struct for storing one emote, with a single from/to position.
// Storing each emote occurance in one object allows us to properly sort the emotes
// to ease the
type Emote struct {
	ID     string `json:"id"`
	From   int    `json:"from"`
	To     int    `json:"to"`
	Source string `json:"source"` // Source is used to allow parsing and inserting more emotes, ex. from BTTV
}

//ByPos is a sorting interface for Emotes
type ByPos []*Emote

func (a ByPos) Len() int           { return len(a) }
func (a ByPos) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPos) Less(i, j int) bool { return a[i].From < a[j].From }

// ParseMessage parses a message from a raw string into a *Message
func ParseMessage(raw string) *Message {
	raw = strings.TrimSpace(raw)
	if len(raw) < 1 {
		return nil
	}
	m := new(Message)

	// Next delimiter, before the next part we want to parse (i)
	// nextDelimiter is always relative to the current position, so actual index
	// is position + nextDelimiter
	nextDelimiter := 0
	// working position/cursor (c)
	position := 0

	//Extract tags
	if raw[0] == prefixTags {
		// If tags are indicated, but no space after (no command), return nil
		if nextDelimiter = strings.IndexByte(raw, space); nextDelimiter < 0 {
			return nil
		}
		m.Tags = ParseTags(raw[1:nextDelimiter])
		position += nextDelimiter + 1
	}

	// Extract and sipmlify prefix as "From";
	// since all twitch users have generic host/nick (user!user@user.tmi.twitch.tv)
	// We only really need a single "from", instead of breaking it up to save arbitrary data
	// Only possible use case would be if we explicitly need to know if a message
	// was sent from a user or if it was sent from server
	if raw[position] == prefix {
		// If prefix is indicated but no space after (command is missing), return nil
		if nextDelimiter = strings.IndexByte(raw[position:], space); position+nextDelimiter < 0 {
			return nil
		}

		m.From = raw[position+1 : position+nextDelimiter]

		if a := strings.IndexByte(m.From, prefixUser); a != -1 {
			m.From = m.From[0:a]
		}
		// Move position forward
		position += nextDelimiter + 1
	}

	// Find end of command
	nextDelimiter = strings.IndexByte(raw[position:], space)
	if nextDelimiter < 0 {
		// Nothing after command, return
		m.Command = raw[position:]
		return m
	}
	m.Command = raw[position : position+nextDelimiter]
	position += nextDelimiter + 1

	// Find prefix for trailing
	nextDelimiter = strings.IndexByte(raw[position:], prefix)

	var params []string
	if nextDelimiter < 0 {
		//no trailing
		params = strings.Split(raw[position:], string(space))
	} else {
		// Has trailing
		if nextDelimiter > 0 {
			// Has params
			params = strings.Split(raw[position:position+nextDelimiter-1], string(space))
		}
		m.Trailing = raw[position+nextDelimiter+1:]
	}
	if len(params) > 0 {
		m.Params = params
	}
	return cleanMessage(m)
}

// Do some default normalising and cleaning
func cleanMessage(m *Message) *Message {
	if m.Command == "PRIVMSG" {
		//Normalise ACTION(/me) on PRIVMSGs
		tLen := len(m.Trailing)
		if tLen > actionPrefixLen {
			if (m.Trailing[0:actionPrefixLen] == actionPrefix) && (m.Trailing[tLen-1] == actionSuffix) {
				m.Trailing = m.Trailing[actionPrefixLen : tLen-1]
				m.Action = true
			}
		}
	}
	return m
}

// ParseEmotes is a short way to automatically parse and save the message's emotes, using tmi.ParseEmotes
func (m *Message) ParseEmotes() {
	if m == nil {
		return
	}
	if m.Tags != nil {
		s, ok := m.Tags["emotes"]
		if !ok {
			return
		}
		m.Emotes = ParseEmotes(s)
		sort.Sort(ByPos(m.Emotes))
	}
}

// Bytes is used to return a Message to a []byte, in case we want to send a *Message to the server
// This does not return a parsed Message to its original form, but rather a message
// in the basic form that the server expects
func (m *Message) Bytes() []byte {
	var buf bytes.Buffer

	buf.WriteString(m.Command)

	if len(m.Params) > 0 {
		buf.WriteByte(space)
		buf.WriteString(strings.Join(m.Params, string(space)))
		buf.WriteString(m.Trailing)
	}
	if buf.Len() > (maxLength) {
		buf.Truncate(maxLength)
	}
	return buf.Bytes()
}

// String returns a stringified version of the message, see Message.Bytes
func (m *Message) String() string {
	return string(m.Bytes())
}

// Channel is a simple method to get the channel, aka the first param
func (m *Message) Channel() string {
	if len(m.Params) > 0 {
		if m.Params[0][0] == '#' {
			return m.Params[0]
		}
	}
	return ""
}

// ParseEmotes transform the emotes string from the tag and returns a slice
// containing individual *Emote instances for each emote occurance.
func ParseEmotes(emoteString string) []*Emote {
	emotes := []*Emote{}
	if emoteString == "" {
		return nil
	}
	emoteSplit := strings.Split(emoteString, "/")
	var occuranceSplit []string
	var split []string
	var id string
	for _, e := range emoteSplit {
		split = strings.Split(e, ":")
		id = split[0]
		occuranceSplit = strings.Split(split[1], ",")

		for _, o := range occuranceSplit {
			i := strings.IndexByte(o, emSep)
			from, _ := strconv.Atoi(o[:i])
			to, _ := strconv.Atoi(o[i+1:])
			emotes = append(emotes, &Emote{ID: id, From: from, To: to, Source: "twitch"})
		}
	}
	return emotes
}

// ParseTags turns the tag prefix string into a proper map[string]string
func ParseTags(s string) map[string]string {
	result := make(map[string]string)

	c, i, i2 := 0, 0, 0 // text cursor and indexes
	for {
		i = strings.IndexByte(s[c:], tagSep) // i = end of command, i2 = equal-sign
		if i > 0 {
			i2 = strings.IndexByte(s[c:c+i], tagAss)
			result[s[c:c+i2]] = s[c+i2+1 : c+i]
			c += i + 1
		} else {
			i2 = strings.IndexByte(s[c:], tagAss)
			result[s[c:c+i2]] = s[c+i2+1:]
			break
		}
	}
	return result
}
