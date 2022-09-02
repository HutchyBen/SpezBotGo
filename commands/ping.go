package commands

import (
	"spezbot/bot"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands = append(Commands, &bot.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "echo",
			Description: "spez is a lil copycat...",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tings",
					Description: "WHAT THE FUCK WILL SPEZ SAY",
					Required:    true,
				},
			},
		},

		Run: func(b *bot.Bot, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
			return &discordgo.InteractionResponseData{
				Content: evt.ApplicationCommandData().Options[0].StringValue(),
			}
		},
	})
}
