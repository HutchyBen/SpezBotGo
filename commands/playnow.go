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
			Name:        "playnow",
			Description: "I AHTE THE CURRENT SONG FUCK OFF NEEK",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "toone",
					Description: "replacion of the now song",
					Required:    true,
				},
			},
		},
		Run: func(b *bot.Bot, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
			songs, vi, isURL, err := PlayStart(b, evt)
			if err != nil {
				return &discordgo.InteractionResponseData{
					Content: "Error: " + err.Error(),
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
					vi.QueueSongNext(evt.Member, songs.Tracks[i])
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
			vi.Guild.Stop()
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
