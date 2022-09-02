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
			Name:        "playnext",
			Description: "vip queue bypass SUCKY SUCKY NIPPIES",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "choon",
					Description: "vascular bypass",
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
			// If no song is playing play the song
			if vi.NowPlaying == nil {
				vi.Guild.PlayTrack(songs.Tracks[0])
				vi.NowPlaying = &bot.NowPlaying{
					Member: evt.Member,
				}
				return &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "Joined Voice Chat",
							Color: 0x00ff00,
						},
					},
				}
			}
			if isURL && len(songs.Tracks) > 1 {
				for i := 0; i < len(songs.Tracks); i++ {
					vi.QueueSong(evt.Member, songs.Tracks[i])
				}
				return &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: fmt.Sprintf("Added %v songs to the queue", len(songs.Tracks)),
							Color: 0x00ff00,
						},
					},
				}
			}
			vi.QueueSongNext(evt.Member, songs.Tracks[0])
			return &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: fmt.Sprintf("Added %v to the queue", songs.Tracks[0].Info.Title),
						Color: 0x00ff00,
					},
				},
			}

		},
	})
}
