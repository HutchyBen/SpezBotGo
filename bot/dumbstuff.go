package bot

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) InitDumbStuff() {
	go b.ImACuteCatboy()
}

func (b *Bot) ImACuteCatboy() {
	for {
		if rand.Intn(99) < 1 {
			g, err := b.Client.State.Guild("802660679856160818")
			if err != nil {
				continue
			}
			c := g.Channels[rand.Intn(len(g.Channels))]
			for c.ID == "1023738749835804672" || c.ID == "1040049481216970794" || c.ID == "802697545893412895" || c.ID == "1003007871916712078" {
				c = g.Channels[rand.Intn(len(g.Channels))]
				// send the message
				_, err = b.Client.ChannelMessageSend(c.ID, "im a cute catboy nyaa~")
				if err != nil {
					continue
				}
			}
		}
		time.Sleep(time.Hour)
	}
}

func (b *Bot) E621Censor(s *discordgo.Session, evt *discordgo.MessageCreate) {
	if strings.Contains(strings.ToLower(evt.Content), "e621") {
		s.ChannelMessageSend(evt.ChannelID, "<@!"+evt.Author.ID+"> Kill yourself freak")
		s.ChannelMessageDelete(evt.ChannelID, evt.ID)
	}
}
