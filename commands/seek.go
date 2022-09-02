package commands

import (
	"spezbot/bot"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands = append(Commands, &bot.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "seek",
			Description: "ok so the first part of the song is bad but the second part is liek really good",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "time",
					Description: "12h11m10s format plz xoxo",
					Required:    true,
				},
			},
		},
		Run: func(b *bot.Bot, ic *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
			vi, ok := b.VoiceInstances[ic.GuildID]
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

			dur, err := time.ParseDuration(ic.ApplicationCommandData().Options[0].StringValue())
			if err != nil {
				return &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Invalid time format (Use 12h11m10s for your time xoxo)",
							Color: 0xff0000,
						},
					},
				}
			}

			err = vi.Guild.Seek(dur)
			if err != nil {
				return &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Error seeking",
							Color: 0xff0000,
						},
					},
				}
			}

			return &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "Seeked to " + dur.String(),
						Color: 0x00ff00,
					},
				},
			}
		},
	})
}
