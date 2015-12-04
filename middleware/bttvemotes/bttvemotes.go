package bttvemotes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/SunspotsEU/tmi"
)

var space = regexp.MustCompile(`\s`)

// BTTVEmote is unmarshaled from the BTTV API
type BTTVEmote struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	Channel   string `json:"channel"`
	Regex     *regexp.Regexp
	ImageType string `json:"imageType"`
}

// BTTVEmoteSet is used to unmarshal sets from the BTTV API
type BTTVEmoteSet struct {
	Status      int          `json:"status"`
	Emotes      []*BTTVEmote `json:"emotes"`
	URLTemplate string       `json:"urlTemplate"`
}

// BTTVEmotes is also used for unmarshalling sets from the BTTV API
type BTTVEmotes struct {
	Sets map[string]*BTTVEmoteSet
}

func (bttv *BTTVEmotes) Download(url string, setName string) {
	// download all emotes

	res, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching BTTV Emotes:", err)
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Println("Error fetching BTTV Emotes:", err)
		return
	}
	var response BTTVEmoteSet
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Error Unmarshaling BTTV Emotes:", err)
		return
	}
	if response.Status != 200 {
		log.Println("Error fetching BTTV Emotes, status:", response.Status)
		return
	}
	log.Println("Fetched", len(response.Emotes), "bttv emotes for", setName)
	bttv.Sets[setName] = &response
}
func (bttv *BTTVEmotes) DownloadChannelEmotes(channel string) {
	bttv.Download("https://api.betterttv.net/2/channels/"+channel, channel)
	bttv.MakeEmoteRegexps(bttv.Sets[channel])
}

func (bttv *BTTVEmotes) DownloadEmotes() {
	bttv.Download("https://api.betterttv.net/2/emotes", "bttv")
	bttv.MakeEmoteRegexps(bttv.Sets["bttv"])
}

func (bttv *BTTVEmotes) MakeEmoteRegexps(emoteRes *BTTVEmoteSet) {
	if emoteRes == nil {
		return
	}
	emotes := emoteRes.Emotes

	for _, emote := range emotes {
		code := regexp.QuoteMeta(emote.Code)
		// `(?:\s|^)` + code + `(?:\s|$)`
		regex, err := regexp.Compile(`(^|\s)` + code + `($|\s)`)
		if err != nil {
			log.Println("Failed to compile emote regex for", emote.Code, " remember me to fix this //Sunspots")
			continue
		}
		emote.Regex = regex
	}
}

func (bttv *BTTVEmotes) MatchEmotes(m *tmi.Message) []*tmi.Emote {
	foundEmotes := []*tmi.Emote{}

	for set, emoteRes := range bttv.Sets {
		for _, emote := range emoteRes.Emotes {
			if emote.Regex == nil {
				continue
			}
			if emote.Channel != "" {
				if strings.ToLower(emote.Channel) != m.Params[0][1:] {
					if set != m.Params[0][1:] {
						continue
					}
				}
			}
			found := emote.Regex.FindAllStringIndex(m.Trailing, -1)
			if found == nil {
				continue
			}
			for _, pos := range found {
				if space.Match([]byte(m.Trailing[pos[0] : pos[0]+1])) {
					pos[0]++
				}
				if pos[1] != len(m.Trailing) {
					if space.MatchString(m.Trailing[pos[1]-1 : pos[1]]) {
						pos[1]--
					}
				}
				foundEmotes = append(foundEmotes, &tmi.Emote{
					ID:     emote.ID,
					From:   pos[0],
					To:     pos[1],
					Source: "bttv",
				})
			}
		}
	}
	return foundEmotes
}

// TBD
func (bttv *BTTVEmotes) MiddleWare(m *tmi.Message, err error) (*tmi.Message, error) {
	if err != nil {
		return m, err
	}
	return m, nil
}

func New() *BTTVEmotes {
	return &BTTVEmotes{Sets: map[string]*BTTVEmoteSet{}}
}
