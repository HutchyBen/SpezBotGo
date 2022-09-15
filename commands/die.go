package commands

import (
	"spezbot/bot"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands = append(Commands, &bot.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "kill",
			Description: "wtf you would kill spez would you...",
		},

		Run: func(b *bot.Bot, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
			vc, ok := b.VoiceInstances[evt.GuildID]
			if !ok {
				return &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Error",
							Description: "I'm already dead SO SAD>>>",
							Color:       0xff0000,
						},
					},
				}
			}

			vc.Suicide()
			return &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Success",
						Description: "I'm dead wtf....",
						Color:       0x00ff00,
					},
				},
			}
		},
	})
}
