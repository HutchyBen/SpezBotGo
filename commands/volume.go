package commands

import (
	"spezbot/bot"

	"github.com/bwmarrin/discordgo"
)

func getFloatPTR(f float64) *float64 {
	return &f
}

func init() {
	Commands = append(Commands, &bot.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "volume",
			Description: "me when erm man or muppet",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "loundyness",
					Description: "HOW FUFFING LOUD",
					Required:    true,
					MinValue:    getFloatPTR(0),
					MaxValue:    1000,
				},
			},
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
			vi.Guild.UpdateVolume(uint16(evt.Interaction.ApplicationCommandData().Options[0].IntValue()))
			return &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: "<:ujel:995463759021748304>",
						Color: 0xffffff,
					},
				},
			}
		},
	})
}
