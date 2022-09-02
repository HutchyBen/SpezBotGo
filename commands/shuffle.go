package commands

import (
	"spezbot/bot"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands = append(Commands, &bot.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "shuffle",
			Description: "whizzle up your queue a lil",
		},

		Run: func(b *bot.Bot, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
			noSongs := &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "You dont have any songs queued",
						Color: 0xffff00,
					},
				},
			}

			vi, ok := b.VoiceInstances[evt.GuildID]
			if !ok {
				return noSongs
			}

			if len(vi.Queues) == 0 {
				return noSongs
			}

			for _, q := range vi.Queues {
				if q.Member.User.ID == evt.Member.User.ID && len(q.Tracks) > 0 {
					q.Shuffle()
					return &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Title: "Shuffling your queue",
								Color: 0x00ff00,
							},
						},
					}
				}
			}
			return noSongs
		},
	})
}
