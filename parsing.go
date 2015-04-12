package tmi

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	actionStart = "\x01ACTION " // Indicates the start of an ACTION(/me) message
	actionEnd   = '\x01'        // Indicates the end of the same ^
)

func cleanEvent(event *Event) *Event {
	// Add "broadcaster" UserType, makes sure that
	if to := event.Channel(); len(to) > 0 && event.From == to[1:] {
		if event.Tags != nil {
			event.Tags.UserType = "broadcaster"
		}
	}
	//Normalise ACTION(/me) on PRIVMSGs
	if event.Command == "PRIVMSG" {
		startLen := len(actionStart)
		if msg := event.Message(); len(msg) > startLen {
			if (msg[0:startLen] == actionStart) && (msg[len(msg)-1] == actionEnd) {
				event.Params[1] = event.Params[1][startLen : len(event.Params[1])-1]
				event.Action = true
			}
		}
	}
	return event
}

func parseEmotes(emoteString string) []*Emote {
	emotes := []*Emote{}
	if emoteString == "" {
		return emotes
	}
	emoteSplit := strings.Split(emoteString, "/")
	for _, e := range emoteSplit {
		split := strings.Split(e, ":")
		id := split[0]
		occuranceSplit := strings.Split(split[1], ",")
		for _, o := range occuranceSplit {
			oSplit := strings.Split(o, "-")
			from, _ := strconv.Atoi(oSplit[0])
			to, _ := strconv.Atoi(oSplit[1])
			emotes = append(emotes, &Emote{Id: id, From: from, To: int(to)})
		}
	}
	return emotes
}

func parseEmoteSets(emoteSetString string) []int {
	if emoteSetString == "" {
		return nil
	}
	stringArr := strings.Split(emoteSetString, ",")
	sets := make([]int, len(stringArr))
	for i, set := range stringArr {
		num, _ := strconv.Atoi(set)
		sets[i] = num
	}
	return sets
}

func makeUserTags(tagMap map[string]string) *UserTags {
	emotes := parseEmotes(tagMap["emotes"])
	sort.Sort(ByPos(emotes))
	emotesets := parseEmoteSets(tagMap["emotesets"])
	tags := &UserTags{
		Color:      tagMap["color"],
		Emotes:     emotes,
		EmoteSets:  emotesets,
		Subscriber: tagMap["subscriber"] == "1",
		Turbo:      tagMap["turbo"] == "1",
		UserType:   tagMap["user_type"],
	}
	return tags
}

func parseTags(tagString string) map[string]string {
	result := make(map[string]string)
	splitTags := strings.Split(tagString, ";")
	for _, tag := range splitTags {
		splitTag := strings.SplitN(tag, "=", 2)
		if len(splitTag) != 2 {
			continue
		}
		result[splitTag[0]] = splitTag[1]
	}
	return result
}

func parseEvent(msg string) *Event {
	msg = strings.TrimSpace(msg)
	if len(msg) < 1 {
		return nil
	}
	event := &Event{Time: time.Now()}

	//Extract tags
	if msg[0] == '@' {
		split := strings.SplitN(msg, " ", 2)
		event.RawTags = parseTags(split[0][1:])
		event.Tags = makeUserTags(event.RawTags)
		msg = split[1]
	}

	//Extract sender/"from"
	//Since all users have the same host, and each user's username and nick is the same,
	//We only really need a single "from", instead of breaking up the prefix
	if msg[0] == ':' {
		split := strings.SplitN(msg, " ", 2)
		event.From = split[0][1:]

		if i := strings.IndexRune(event.From, '!'); i != -1 {
			event.From = event.From[0:i]
		}

		msg = split[1]
	}

	// trailing = split[1]
	split := strings.SplitN(msg, " :", 2)
	params := strings.Split(split[0], " ")
	event.Command = strings.ToUpper(params[0])
	event.Params = params[1:]
	if len(split) > 1 {
		event.Params = append(event.Params, split[1])
	}

	return cleanEvent(event)
}
