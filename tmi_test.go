package tmi

import (
	"log"
	"testing"
)

// I'm pretty bad at making tests, but I made some pretty basic ones anyway, to start somewhere.

func TestParseEvent(t *testing.T) {
	//There's probably some proper eval stuff available, but I made this to do quick string comparisons.
	evalString := func(expected string, result string) {
		if expected != result {
			t.Errorf("expected  %s, got %v", expected, result)
		}
	}
	var evt *Event

	// Empty string
	evt = parseEvent("")
	if evt != nil {
		t.Errorf("Empty string did not return nil", evt)
	}

	// Command-only message
	// I don't think this is right, but we'll let it pass
	evt = parseEvent("RECONNECT")
	if evt == nil {
		t.Errorf("Event was nil")
	}
	evalString(evt.Message(), "")
	evalString(evt.Channel(), "")
	evalString(evt.From, "")

	// Test command with channel target
	evt = parseEvent(":sunsbot!sunsbot@sunsbot.tmi.twitch.tv JOIN #sunspots")
	evalString("sunsbot", evt.From)
	evalString("JOIN", evt.Command)
	evalString("#sunspots", evt.Message()) //Should this be the message? Why not?
	evalString("#sunspots", evt.Channel())

	// Test command without prefix
	evt = parseEvent("PING :tmi.twitch.tv")
	evalString("PING", evt.Command)
	evalString("tmi.twitch.tv", evt.Message())
	evalString("", evt.Channel()) //Make sure it doesn't return anything as channel

	// Test channel PRIVMSG without tags
	evt = parseEvent(":sunspots!sunspots@sunspots.tmi.twitch.tv PRIVMSG #sunspots :test message, lol")
	evalString("#sunspots", evt.Channel())
	evalString("test message, lol", evt.Message())
	evalString("sunspots", evt.From)

	// Test command without trailing(message)
	evt = parseEvent(":jtv MODE #sunspots +o sunspots")
	evalString("MODE", evt.Command)
	evalString("#sunspots", evt.Channel())
	evalString("+o", evt.Params[1])
	evalString("sunspots", evt.Params[2])

	// Test PRIVMSG with tags
	evt = parseEvent("@color=#FF6BFF;emotes=;subscriber=0;turbo=0;user_type= :sunspots!sunspots@sunspots.tmi.twitch.tv PRIVMSG #sunspots :test message, lol")
	evalString("sunspots", evt.From)
	evalString("#FF6BFF", evt.Tags.Color)
	evalString("broadcaster", evt.Tags.UserType) // Sender is same as channel
	if evt.Tags.Subscriber != false {
		t.Errorf("expected  %s, got %v", false, evt.Tags.Subscriber)
	}
	if evt.Tags.Turbo != false {
		t.Errorf("expected  %s, got %v", false, evt.Tags.Turbo)
	}

	//Test PRIVMSG with ACTION(/me) and tags
	evt = parseEvent("@color=#FF6BFF;emotes=;subscriber=0;turbo=0;user_type= :sunspots!sunspots@sunspots.tmi.twitch.tv PRIVMSG #sunspots :ACTION likes pie")
	evalString(evt.Message(), "likes pie")
	if !evt.Action {
		t.Errorf("Action was not true for action message")
	}

	evt = parseEvent("@color=#FF6BFF;emotes=15762:0-4,12-16/1902:6-10;subscriber=0;turbo=0;user_type= :sunspots!sunspots@sunspots.tmi.twitch.tv PRIVMSG #sunspots :Kappa Keepo Kappa")
	//Find a good way to evaluate the emotes...
	//log.Println(evt)
}

func TestParseTags(t *testing.T) {
	evalString := func(expected string, result string) {
		if expected != result {
			t.Errorf("expected  %s, got %v", expected, result)
		}
	}
	var tags map[string]string
	tags = parseTags("color=#FF6BFF;emotes=;subscriber=0;turbo=0;user_type")
	log.Println(tags)
	evalString("#FF6BFF", tags["color"])
	evalString("", tags["emotes"])
	evalString("0", tags["subscriber"])
	evalString("0", tags["turbo"])
	evalString("", tags["user_type"])
}
