package tmi

import (
	"strconv"
	"strings"
)

func parseEmotes(emoteString string) []*Emote {
	emotes := []*Emote{}
	if emoteString == "" {
		return nil
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
			emotes = append(emotes, &Emote{ID: id, From: from, To: int(to)})
		}
	}
	return emotes
}
