package tmi

import (
	"reflect"
	"testing"
)

var testMessages = map[string]*Message{
	// Message with prefix, command and one argument
	":sunsbot!sunsbot@sunsbot.tmi.twitch.tv JOIN #sunspots": &Message{
		Raw:     ":sunsbot!sunsbot@sunsbot.tmi.twitch.tv JOIN #sunspots",
		From:    "sunsbot",
		Command: "JOIN",
		Params:  []string{"#sunspots"},
		Channel: "sunspots",
	},
	// PRIVMSG message
	":sunspots!sunspots@sunspots.tmi.twitch.tv PRIVMSG #sunspots :test message, lol": &Message{
		Raw:     ":sunspots!sunspots@sunspots.tmi.twitch.tv PRIVMSG #sunspots :test message, lol",
		From:    "sunspots",
		Command: "PRIVMSG",
		Params:  []string{"#sunspots"},
		Channel: "sunspots",
		Text:    "test message, lol",
	},
	// Message with tags
	"@color=#FF6BFF;emotes=;subscriber=0;turbo=0;user-type= :sunsbot!sunsbot@sunsbot.tmi.twitch.tv PRIVMSG #sunspots :test message, lol": &Message{
		Raw:     "@color=#FF6BFF;emotes=;subscriber=0;turbo=0;user-type= :sunsbot!sunsbot@sunsbot.tmi.twitch.tv PRIVMSG #sunspots :test message, lol",
		From:    "sunsbot",
		Command: "PRIVMSG",
		Params:  []string{"#sunspots"},
		Channel: "sunspots",
		Text:    "test message, lol",
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
		Raw:     "@color=#FF6BFF;emotes=;subscriber=0;turbo=0;user-type= :sunspots!sunspots@sunspots.tmi.twitch.tv PRIVMSG #sunspots :test message, lol",
		From:    "sunspots",
		Command: "PRIVMSG",
		Params:  []string{"#sunspots"},
		Channel: "sunspots",
		Text:    "test message, lol",
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
		Raw:     "@color=#FF6BFF;emotes=;subscriber=0;turbo=0;user-type= :sunspots!sunspots@sunspots.tmi.twitch.tv PRIVMSG #sunspots :\001ACTION likes pie\001",
		From:    "sunspots",
		Command: "PRIVMSG",
		Params:  []string{"#sunspots"},
		Channel: "sunspots",
		Text:    "likes pie",
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
		Raw:     "@color=#FF6BFF;display-name=Sunspots;emotes=25:0-4/1902:12-16/30259:23-29,31-37;subscriber=0;turbo=0;user-type= :sunspots!sunspots@sunspots.tmi.twitch.tv PRIVMSG #sunspots :Kappa Hello Keepo test HeyGuys HeyGuys",
		From:    "sunspots",
		Command: "PRIVMSG",
		Params:  []string{"#sunspots"},
		Channel: "sunspots",
		Text:    "Kappa Hello Keepo test HeyGuys HeyGuys",
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
	"@color=#8A2BE2;display-name=Someperson\\s\\s;emotes=;subscriber=0;turbo=0;user-type= :someperson!someperson@someperson.tmi.twitch.tv PRIVMSG #sunspots :Hello!": &Message{
		Raw:     "@color=#8A2BE2;display-name=Someperson\\s\\s;emotes=;subscriber=0;turbo=0;user-type= :someperson!someperson@someperson.tmi.twitch.tv PRIVMSG #sunspots :Hello!",
		From:    "someperson",
		Command: "PRIVMSG",
		Params:  []string{"#sunspots"},
		Channel: "sunspots",
		Text:    "Hello!",
		Tags: map[string]string{
			"display-name": "Someperson\\s\\s",
			"color":        "#8A2BE2",
			"emotes":       "",
			"subscriber":   "0",
			"turbo":        "0",
			"user-type":    "",
		},
	},
	// Message with unicode
	"@badges=subscriber/0;color=;display-name=좋은녀석;emotes=;flags=;id=a227f2e2-23c1-320c-adaf-f6fb877c04ae;mod=0;room-id=678912345;subscriber=1;tmi-sent-ts=1551880839830;turbo=0;user-id=123456789;user-type= :user123!user123@user123.tmi.twitch.tv PRIVMSG #somechannel :ㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑ": &Message{
		Raw:     "@badges=subscriber/0;color=;display-name=좋은녀석;emotes=;flags=;id=a227f2e2-23c1-320c-adaf-f6fb877c04ae;mod=0;room-id=678912345;subscriber=1;tmi-sent-ts=1551880839830;turbo=0;user-id=123456789;user-type= :user123!user123@user123.tmi.twitch.tv PRIVMSG #somechannel :ㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑ",
		From:    "user123",
		Command: "PRIVMSG",
		Params:  []string{"#somechannel"},
		Channel: "somechannel",
		Text:    "ㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑㅑ",
		Tags: map[string]string{
			"badges":       "subscriber/0",
			"color":        "",
			"display-name": "좋은녀석",
			"emotes":       "",
			"flags":        "",
			"id":           "a227f2e2-23c1-320c-adaf-f6fb877c04ae",
			"mod":          "0",
			"room-id":      "678912345",
			"subscriber":   "1",
			"tmi-sent-ts":  "1551880839830",
			"turbo":        "0",
			"user-id":      "123456789",
			"user-type":    "",
		},
	},
	// Message with no params but with trailing
	":sunspots!sunspots@sunspots.tmi.twitch.tv QUIT :Went somewhere else": &Message{
		Raw:     ":sunspots!sunspots@sunspots.tmi.twitch.tv QUIT :Went somewhere else",
		From:    "sunspots",
		Command: "QUIT",
		Text:    "Went somewhere else",
	},
	// Message with prefix and command, no trailing
	":tmi.twitch.tv RECONNECT": &Message{
		Raw:     ":tmi.twitch.tv RECONNECT",
		Command: "RECONNECT",
		From:    "tmi.twitch.tv",
	},
	// Message with command only
	"RECONNECT": &Message{
		Raw:     "RECONNECT",
		Command: "RECONNECT",
	},
	// Message with command and trailing, no prefix
	"PRIVMSG :hello there": &Message{
		Raw:     "PRIVMSG :hello there",
		Command: "PRIVMSG",
		Text:    "hello there",
	},
	// Catch correct but meaningless stuff, make sure nothing breaks.
	"arbitrary random words": &Message{
		Raw:     "arbitrary random words",
		Command: "arbitrary",
		Params:  []string{"random", "words"},
	},
	"a b c": &Message{
		Raw:     "a b c",
		Command: "a",
		Params:  []string{"b", "c"},
	},
	// Since we just do a simple TrimSpace, it trims arbitrary leading spaces
	" 	leading space": &Message{
		Raw: " 	leading space",
		Command: "leading",
		Params:  []string{"space"},
	},
	"\\s": &Message{
		Raw:     "\\s",
		Command: "\\s",
	},
}

var badMessages = []string{
	// Catch malformed messages, these should not be parseable.
	":Kappa",
	":",
	"\r\n",
	"",
	" ",
	"\n",
}

func logUnEqual(t *testing.T, got, exp interface{}) {
	if got != exp {
		t.Logf("Got '%s', expected '%s'", got, exp)
	}
}

func logDeepUnEqual(t *testing.T, got, exp interface{}) {
	if !reflect.DeepEqual(got, exp) {
		t.Logf("Got: '%s' \n Expected: '%s'", got, exp)
	}
}

func TestParseMessage(t *testing.T) {
	for raw, p := range testMessages {
		m := Message{}
		parseMessage(raw, &m)
		m.ParseEmotes()
		if !reflect.DeepEqual(&m, p) {
			t.Errorf("Failed parsing message: \"%s\"", raw)
			t.Log(m, *p)
			// Try to log the offending message part
			t.Log("Bad actors:")
			logUnEqual(t, m.Raw, p.Raw)
			logUnEqual(t, m.Action, p.Action)
			logUnEqual(t, m.Command, p.Command)
			logUnEqual(t, m.From, p.From)
			logDeepUnEqual(t, m.Params, p.Params)
			logUnEqual(t, m.Channel, p.Channel)
			logDeepUnEqual(t, m.Tags, p.Tags)
			logDeepUnEqual(t, m.Emotes, p.Emotes)
		}
	}
}

func TestBadMessages(t *testing.T) {
	for _, raw := range badMessages {
		m := Message{}
		err := parseMessage(raw, &m)
		if err == nil {
			t.Error("Message parsing should have failed", m)
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

func BenchmarkParseMessage(b *testing.B) {
	out := Message{}
	for i := 0; i < b.N; i++ {
		parseMessage(
			"@badges=;color=#FF0000;display-name=user123;emotes=;flags=;id=3c44363a-c56a-4d42-a040-b939c6fb58f1;mod=0;room-id=56781234;subscriber=0;tmi-sent-ts=1546212821111;turbo=0;user-id=12234234;user-type= :user123!user123@user123.tmi.twitch.tv PRIVMSG #summit1g :big head no good",
			&out)
	}
}

func BenchmarkTags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tags := map[string]string{}
		parseTags("badges=;color=#FF0000;display-name=user123;emotes=;flags=;id=3c44363a-c56a-4d42-a040-b939c6fb58f1;mod=0;room-id=56781234;subscriber=0;tmi-sent-ts=1546212821111;turbo=0;user-id=12234234;user-type=", tags)
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

func BenchmarkCleanMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := &Message{
			Command: "PRIVMSG",
			Text:    "\001ACTION likes pie\001",
		}
		cleanMessage(m)
	}
}
