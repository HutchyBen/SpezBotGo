package commands

import (
	"fmt"
	"spezbot/bot"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands = append(Commands, &bot.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "skip",
			Description: "skip the fucking song xoxoox",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "tings",
					Description: "how many songs to skip",
					Required:    false,
				},
			},
		},

		Run: func(b *bot.Bot, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
			vi, ok := b.VoiceInstances[evt.GuildID]
			if !ok {
				return &discordgo.InteractionResponseData{
					Content: "The bot is not in a voice channel",
				}
			}
			var skipAmount int64 = 1
			if len(evt.ApplicationCommandData().Options) != 0 {
				skipAmount = evt.ApplicationCommandData().Options[0].IntValue()
			}
			var i int64 = 0
			// i have never done for loop like this
			for i = 0; i < skipAmount-1; i++ {
				vi.Queues[vi.QueueIndex].Pop()
				if len(vi.Queues[vi.QueueIndex].Tracks) == 0 {
					vi.Queues = append(vi.Queues[:vi.QueueIndex], vi.Queues[vi.QueueIndex+1:]...)
					vi.QueueIndex--
				}
				vi.QueueIndex++
				if vi.QueueIndex >= len(vi.Queues) {
					vi.QueueIndex = 0
				}

			}
			var stylePoints = " songs"
			if skipAmount == 1 {
				stylePoints = " song"
			}
			vi.Guild.Stop()
			return &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: fmt.Sprintf("Skipped %v%v", skipAmount, stylePoints),
						Color: 0x00ff00,
					},
				},
			}
		},
	})
}
