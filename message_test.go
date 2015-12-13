package tmi

import (
	"reflect"
	"testing"
)

var testMessages = map[string]*Message{
	// Message with prefix, command and one argument
	":sunsbot!sunsbot@sunsbot.tmi.twitch.tv JOIN #sunspots": &Message{
		From:    "sunsbot",
		Command: "JOIN",
		Params:  []string{"#sunspots"},
	},
	// PRIVMSG message
	":sunspots!sunspots@sunspots.tmi.twitch.tv PRIVMSG #sunspots :test message, lol": &Message{
		From:     "sunspots",
		Command:  "PRIVMSG",
		Params:   []string{"#sunspots"},
		Trailing: "test message, lol",
	},
	// Message with tags
	"@color=#FF6BFF;emotes=;subscriber=0;turbo=0;user-type= :sunsbot!sunsbot@sunsbot.tmi.twitch.tv PRIVMSG #sunspots :test message, lol": &Message{
		From:     "sunsbot",
		Command:  "PRIVMSG",
		Params:   []string{"#sunspots"},
		Trailing: "test message, lol",
		Tags: map[string]string{
			"color":      "#FF6BFF",
			"emotes":     "",
			"subscriber": "0",
			"turbo":      "0",
			"user-type":  "",
		},
	},
	// Message with tags, by broadcaster
	"@color=#FF6BFF;emotes=;subscriber=0;turbo=0;user-type= :sunspots!sunspots@sunspots.tmi.twitch.tv PRIVMSG #sunspots :test message, lol": &Message{
		From:     "sunspots",
		Command:  "PRIVMSG",
		Params:   []string{"#sunspots"},
		Trailing: "test message, lol",
		Tags: map[string]string{
			"color":      "#FF6BFF",
			"emotes":     "",
			"subscriber": "0",
			"turbo":      "0",
			"user-type":  "",
		},
	},
	// Message with ACTION
	"@color=#FF6BFF;emotes=;subscriber=0;turbo=0;user-type= :sunspots!sunspots@sunspots.tmi.twitch.tv PRIVMSG #sunspots :\001ACTION likes pie\001": &Message{
		From:     "sunspots",
		Command:  "PRIVMSG",
		Params:   []string{"#sunspots"},
		Trailing: "likes pie",
		Tags: map[string]string{
			"color":      "#FF6BFF",
			"emotes":     "",
			"subscriber": "0",
			"turbo":      "0",
			"user-type":  "",
		},
		Action: true,
	},
	// Message with emotes
	"@color=#FF6BFF;display-name=Sunspots;emotes=25:0-4/1902:12-16/30259:23-29,31-37;subscriber=0;turbo=0;user-type= :sunspots!sunspots@sunspots.tmi.twitch.tv PRIVMSG #sunspots :Kappa Hello Keepo test HeyGuys HeyGuys": &Message{
		From:     "sunspots",
		Command:  "PRIVMSG",
		Params:   []string{"#sunspots"},
		Trailing: "Kappa Hello Keepo test HeyGuys HeyGuys",
		Tags: map[string]string{
			"display-name": "Sunspots",
			"color":        "#FF6BFF",
			"emotes":       "25:0-4/1902:12-16/30259:23-29,31-37",
			"subscriber":   "0",
			"turbo":        "0",
			"user-type":    "",
		},
		Emotes: []*Emote{
			&Emote{ID: "25", From: 0, To: 4, Source: "twitch"},
			&Emote{ID: "1902", From: 12, To: 16, Source: "twitch"},
			&Emote{ID: "30259", From: 23, To: 29, Source: "twitch"},
			&Emote{ID: "30259", From: 31, To: 37, Source: "twitch"},
		},
	},
	// Message with escaped spaces in tags
	"@color=#8A2BE2;display-name=Sparklingsandshrew\\s\\s;emotes=;subscriber=0;turbo=0;user-type= :sparklingsandshrew!sparklingsandshrew@sparklingsandshrew.tmi.twitch.tv PRIVMSG #sunspots :Hello!": &Message{
		From:     "sparklingsandshrew",
		Command:  "PRIVMSG",
		Params:   []string{"#sunspots"},
		Trailing: "Hello!",
		Tags: map[string]string{
			"display-name": "Sparklingsandshrew\\s\\s",
			"color":        "#8A2BE2",
			"emotes":       "",
			"subscriber":   "0",
			"turbo":        "0",
			"user-type":    "",
		},
	},
	// Message with no params but with trailing
	":sunspots!sunspots@sunspots.tmi.twitch.tv QUIT :Went somewhere else": &Message{
		From:     "sunspots",
		Command:  "QUIT",
		Trailing: "Went somewhere else",
	},
	// Message with prefix and command, no trailing
	":tmi.twitch.tv RECONNECT": &Message{
		Command: "RECONNECT",
		From:    "tmi.twitch.tv",
	},
	// Message with command only
	"RECONNECT": &Message{
		Command: "RECONNECT",
	},
	// Message with command and trailing, no prefix
	"PRIVMSG :hello there": &Message{
		Command:  "PRIVMSG",
		Trailing: "hello there",
	},
	// Catch stupid/malformed/empty stuff, make sure nothing breaks.
	"arbitrary bullshit yo": &Message{
		Command: "arbitrary",
		Params:  []string{"bullshit", "yo"},
	},
	"a b c": &Message{
		Command: "a",
		Params:  []string{"b", "c"},
	},
	// Since we just do a simple TrimSpace, it trims arbitrary leading spaces
	" 	leading space": &Message{
		Command: "leading",
		Params:  []string{"space"},
	},
	"\\s": &Message{
		Command: "\\s",
	},
	":poopie": nil,
	":":       nil,
	"\r\n":    nil,
	"":        nil,
	" ":       nil,
	"\n":      nil,
}

func TestParseMessage(t *testing.T) {
	for raw, p := range testMessages {
		m := ParseMessage(raw)
		if m != nil {
			m.ParseEmotes()
		}
		if (p == nil) && (m != nil) {
			t.Errorf("Failed parsing message: \"%s\"", raw)
			t.Log("Expected nil, got", m)
			continue
		}
		if !reflect.DeepEqual(m, p) {
			t.Errorf("Failed parsing message: \"%s\"", raw)
			t.Log(m.Action == p.Action)
			t.Log(m.Command == p.Command)
			t.Log(m.Command, p.Command)
			t.Log(m.From == p.From)
			t.Log(reflect.DeepEqual(m.Params, p.Params))
			t.Log(reflect.DeepEqual(m.Tags, p.Tags))
			t.Log(m.Tags, p.Tags)
			t.Log(reflect.DeepEqual(m.Emotes, p.Emotes))
		}
	}

}
func TestParseEmotes(t *testing.T) {
	emotes := ParseEmotes("25:0-4,6-10,12-16,18-22,24-28,30-34,36-40,42-46,48-52,54-58,60-64,66-70,72-76,78-82,84-88,90-94,96-100,102-106,108-112,114-118,120-124,126-130")
	comp := []*Emote{
		&Emote{ID: "25", From: 0, To: 4, Source: "twitch"},
		&Emote{ID: "25", From: 6, To: 10, Source: "twitch"},
		&Emote{ID: "25", From: 12, To: 16, Source: "twitch"},
		&Emote{ID: "25", From: 18, To: 22, Source: "twitch"},
		&Emote{ID: "25", From: 24, To: 28, Source: "twitch"},
		&Emote{ID: "25", From: 30, To: 34, Source: "twitch"},
		&Emote{ID: "25", From: 36, To: 40, Source: "twitch"},
		&Emote{ID: "25", From: 42, To: 46, Source: "twitch"},
		&Emote{ID: "25", From: 48, To: 52, Source: "twitch"},
		&Emote{ID: "25", From: 54, To: 58, Source: "twitch"},
		&Emote{ID: "25", From: 60, To: 64, Source: "twitch"},
		&Emote{ID: "25", From: 66, To: 70, Source: "twitch"},
		&Emote{ID: "25", From: 72, To: 76, Source: "twitch"},
		&Emote{ID: "25", From: 78, To: 82, Source: "twitch"},
		&Emote{ID: "25", From: 84, To: 88, Source: "twitch"},
		&Emote{ID: "25", From: 90, To: 94, Source: "twitch"},
		&Emote{ID: "25", From: 96, To: 100, Source: "twitch"},
		&Emote{ID: "25", From: 102, To: 106, Source: "twitch"},
		&Emote{ID: "25", From: 108, To: 112, Source: "twitch"},
		&Emote{ID: "25", From: 114, To: 118, Source: "twitch"},
		&Emote{ID: "25", From: 120, To: 124, Source: "twitch"},
		&Emote{ID: "25", From: 126, To: 130, Source: "twitch"},
	}

	if !reflect.DeepEqual(emotes, comp) {
		t.Error("Failed parsing emotes!")
	}
}

func BenchmarkMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseMessage("@color=#FF6BFF;emotes=;subscriber=0;turbo=0;user-type= :sunspots!sunspots@sunspots.tmi.twitch.tv PRIVMSG #sunspots :test message, lol")
	}
}

func BenchmarkTags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseTags("color=#FF6BFF;emotes=;subscriber=0;turbo=0;user-type=")
	}
}

func BenchmarkEmotes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseEmotes("25:0-4/1902:12-16/30259:23-29,31-37")
	}
}

func BenchmarkEmotes22Kappas(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseEmotes("25:0-4,6-10,12-16,18-22,24-28,30-34,36-40,42-46,48-52,54-58,60-64,66-70,72-76,78-82,84-88,90-94,96-100,102-106,108-112,114-118,120-124,126-130")
	}
}
