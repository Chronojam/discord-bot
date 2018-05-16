// Credits to Tristan Colgate-McFarlane, who wrote the original file here:
// https://github.com/tcolgate/hugot/blob/bab4d700f10a69412f85315b66c6e562691b83af/handlers/hears/tableflip/tableflip.go

package tableflip

import (
	bot "github.com/chronojam/discord-bot/pkg/discord-bot"
	"github.com/bwmarrin/discordgo"
	"time"
	"regexp"
	"strings"
)


func init() {
	bot.AddHandler(messageCreate)
}

// We'll be horrid and use some globals
var flipState bool
var lastFlip time.Time
var tableflipRegexp = regexp.MustCompile(`(^| *)tableflip($| *)`)

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {	
	ok, args := bot.ParseMessage(m)
	if !ok {
		return
	}

	if args[1] != "tableflip" {
		return
	}

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	flip := `(╯°□°）╯︵ ┻━┻`
	unFlip := `┬━┬ ノ( ゜-゜ノ)`
	doubleFlip := "┻━┻ ︵¯\\(ツ)/¯ ︵ ┻━┻"
	tripleFlip := "(╯°□°）╯¸.·´¯`·.¸¸.·´¯`·.¸¸.·´¯ ┻━┻"
	flipOff := "ಠ︵ಠ凸"
	
	if !flipState {
		flipState = true
	
		go func() {
			five, _ := time.ParseDuration("30s")
			time.Sleep(five)
			if flipState == true {
				flipState = false
				s.ChannelMessageSend(m.ChannelID, unFlip)
			}
		}()

		txt := strings.Join(args[1:], "")
	
		switch fs := tableflipRegexp.FindAllString(txt, 5); len(fs) {
		case 1:
			s.ChannelMessageSend(m.ChannelID, flip)
		case 2:
			s.ChannelMessageSend(m.ChannelID, doubleFlip)
		case 3:
			s.ChannelMessageSend(m.ChannelID, tripleFlip)
		default:
			s.ChannelMessageSend(m.ChannelID, flipOff)
			flipState = false
		}
		return
	}
	
	flipState = false
	s.ChannelMessageSend(m.ChannelID, unFlip)
}