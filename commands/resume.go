package commands

import (
	"spezbot/bot"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands = append(Commands, &bot.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "resume",
			Description: "Do longer.",
		},

		Run: func(b *bot.Bot, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
			vi, ok := b.VoiceInstances[evt.GuildID]
			if !ok {
				return &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "The bot is not connected to a voice chat",
							Color: 0xff0000,
						},
					},
				}
			}
			vi.Guild.SetPaused(false)

			return &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "Resumed " + vi.NowPlaying.Track.Title,
						Color: 0x00ff00,
					},
				},
			}
		},
	})
}
