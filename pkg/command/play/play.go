package play

import (
	bot "github.com/chronojam/discord-bot/pkg/discord-bot"
	"github.com/bwmarrin/discordgo"
	"encoding/binary"
	"time"
	"io"
	"os"
	"fmt"
	"io/ioutil"
	"net/http"
)


func init() {
	bot.AddHandler(playaSound)
	bot.AddHandler(uploadSound)
}

func uploadSound(s *discordgo.Session, m *discordgo.MessageCreate) {
	ok, args := bot.ParseMessage(m)
	if !ok {
		return
	}

	if args[1] != "upload-sound" {
		return
	}

	for _, a := range m.Attachments {
		//fmt.Println(a.ProxyURL)
		resp, err := http.Get(a.URL)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Write to filesystem.
		err = ioutil.WriteFile("sounds/"+a.Filename, b, 0755)
		if err != nil {
			fmt.Println(err)
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(`Jolly Good. Ill store that one for later.
Play me by calling "!hugot play-sound %s" :musical_score:  
		`, a.Filename))
	}
}

func playaSound(s *discordgo.Session, m *discordgo.MessageCreate) {	
	ok, args := bot.ParseMessage(m)
	if !ok {
		return
	}

	if args[1] != "play-sound" {
		return
	}
	
	// Find the channel that the message came from.
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		// Could not find channel.
		return
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		// Could not find guild.
		return
	}

	// Look for the message sender in that guild's current voice states.
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			err = playSound(s, g.ID, vs.ChannelID, "sounds/" + args[2])
			if err != nil {
				fmt.Println("Error playing sound:", err)
			}

			return
		}
	}
}

func loadSound(path string) ([][]byte, error) {
	var buffer = make([][]byte, 0)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return buffer, err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return buffer, err
			}
			return buffer, nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return buffer, err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return buffer, err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}

	return buffer, nil
}

func playSound(s *discordgo.Session, guildID, channelID, path string) error {
	buffer, err := loadSound(path)
	if err != nil {
		return err
	}
	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}