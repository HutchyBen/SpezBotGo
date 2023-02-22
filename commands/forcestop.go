package commands

import (
	"spezbot/bot"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands = append(Commands, &bot.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "forcestop",
			Description: "so you like rip spezzes spine out or something",
		},

		Run: func(b *bot.Bot, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
			if evt.Member.User.ID == "386912562937331725" {
				b.Die <- true
				return &discordgo.InteractionResponseData{
					Content: "focking dead m8",
				}
			}

			return &discordgo.InteractionResponseData{
				Content: "nah only ben can rip my spine out ",
			}
		},
	})
}
