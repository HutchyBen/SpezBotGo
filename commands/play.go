package commands

import (
	"fmt"
	"spezbot/bot"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands = append(Commands, &bot.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "play",
			Description: "Plays a song or playlist",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "The query to search for",
					Required:    true,
				},
			},
		},
		Run: func(b *bot.Bot, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
			songs, vi, isURL, err := PlayStart(b, evt)
			if err != nil {
				return &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Error",
							Description: strings.ToUpper(err.Error()[:1]) + err.Error()[1:],
							Color:       0xff0000,
						},
					},
				}
			}
			embedTitle := ""

			if isURL && len(songs.Tracks) > 1 {
				for i := 0; i <= len(songs.Tracks)-1; i++ {
					vi.QueueSong(evt.Member, songs.Tracks[i])
				}
				embedTitle = fmt.Sprintf("Added %v songs to the queue", len(songs.Tracks))
			} else {
				vi.QueueSong(evt.Member, songs.Tracks[0])
				embedTitle = fmt.Sprintf("Added %v to the queue", songs.Tracks[0].Info.Title)
			}

			if vi.NowPlaying == nil {
				err = vi.Guild.PlayTrack(*vi.Queues[0].Pop())
				// Loop unitil song that isnt borked is found or until queue is empty
				for err != nil {
					if len(vi.Queues[0].Tracks) > 0 {
						err = vi.Guild.PlayTrack(*vi.Queues[0].Pop())
					} else {
						return &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{
								{
									Title: "Could not play song",
									Color: 0xff0000,
								},
							},
						}
					}
				}
				vi.NowPlaying = &bot.NowPlaying{
					Member: evt.Member,
				}
				// should be empty
				if len(vi.Queues[0].Tracks) == 0 {
					vi.Queues = make([]bot.Queue, 0)
				}
				return &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Starting voice session",
							Color: 0x00ff00,
						},
					},
				}
			}

			return &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: embedTitle,
						Color: 0x00ff00,
					},
				},
			}

		},
	})
}
