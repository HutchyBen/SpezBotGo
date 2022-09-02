package commands

import (
	"fmt"
	"spezbot/bot"
	"time"

	"github.com/bwmarrin/discordgo"
)

func SplitIntoPages(songs []string) []string {
	pages := []string{}

	for _, song := range songs {
		if len(pages) == 0 {
			pages = append(pages, song)
			continue
		}

		if len(pages[len(pages)-1])+len(song) >= 2048 {
			pages = append(pages, song)
			continue
		}

		pages[len(pages)-1] += song
	}
	return pages
}

func MakeEdit(pages []string, pageNum int) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Queue",
					Description: pages[pageNum],
					Color:       0x00ff00,
					Footer: &discordgo.MessageEmbedFooter{
						Text: fmt.Sprintf("Page %v/%v", pageNum+1, len(pages)),
					},
				},
			},
			Components: []discordgo.MessageComponent{
				&discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "◀️",
							Style:    discordgo.SuccessButton,
							CustomID: "previous",
						},
						&discordgo.Button{
							Label:    "▶️",
							Style:    discordgo.SuccessButton,
							CustomID: "next",
						},
					},
				},
			},
		},
	}
}
func init() {
	Commands = append(Commands, &bot.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "queue",
			Description: "spez likes his music ( ͡° ͜ʖ ͡°)",
		},

		Run: func(b *bot.Bot, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
			vi, ok := b.VoiceInstances[evt.GuildID]
			if !ok {
				return &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "The bot is not in a voice channel",
							Color: 0xff0000,
						},
					},
				}
			}

			if len(vi.Queues) == 0 {
				return &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title: "The queue is empty",
							Color: 0xffff00,
						},
					},
				}
			}

			songList := []string{}

			qIdx := vi.QueueIndex + 1
			songPlace := make([]int, len(vi.Queues))
			i := 1
			for {
				if qIdx >= len(vi.Queues) {
					qIdx = 0
				}

				if songPlace[qIdx] >= len(vi.Queues[qIdx].Tracks) {
					qIdx++

					for i := 0; i < len(vi.Queues)-1; i++ {
						if songPlace[i] < len(vi.Queues[i].Tracks) {
							continue
						}
					}
					break
				}

				songList = append(songList, fmt.Sprintf("%v. %v <@%v>\n", i, vi.Queues[qIdx].Tracks[songPlace[qIdx]].Info.Title, vi.Queues[qIdx].Member.User.ID))
				songPlace[qIdx]++
				qIdx++
				i++
			}
			pages := SplitIntoPages(songList)
			pageNum := 0

			b.Client.InteractionRespond(evt.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Queue",
							Description: pages[pageNum],
							Color:       0x00ff00,
							Footer: &discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("Page %v/%v", pageNum+1, len(pages)),
							},
						},
					},
					Components: []discordgo.MessageComponent{
						&discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								&discordgo.Button{
									Label:    "◀️",
									Style:    discordgo.SuccessButton,
									CustomID: "previous",
								},
								&discordgo.Button{
									Label:    "▶️",
									Style:    discordgo.SuccessButton,
									CustomID: "next",
								},
							},
						},
					},
				},
			})

			resp := make(chan *discordgo.InteractionCreate)

			b.CH.Events[evt.GuildID] = resp
			timer := time.NewTimer(time.Minute)

			for {
				select {
				case <-timer.C:
					goto loopend
				case e := <-resp:
					switch e.MessageComponentData().CustomID {
					case "previous":
						if pageNum == 0 {
							pageNum = len(pages) - 1
						} else {
							pageNum--
						}
						b.Client.InteractionRespond(e.Interaction, MakeEdit(pages, pageNum))
					case "next":
						if pageNum == len(pages)-1 {
							pageNum = 0
						} else {
							pageNum++
						}
						b.Client.InteractionRespond(e.Interaction, MakeEdit(pages, pageNum))

					}
				}
			}
		loopend:
			return nil
		},
	})
}
