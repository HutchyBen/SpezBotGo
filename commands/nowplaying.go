package commands

import (
	"fmt"
	"math"
	"spezbot/bot"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func DescProgress(now uint, total uint) string {
	nowTime := time.Duration(now) * time.Millisecond
	totalTime := time.Duration(total) * time.Millisecond
	var prog = float64(nowTime) / float64(totalTime)
	var emoPos = math.Round(10 * prog)
	str := "▬▬▬▬▬▬▬▬▬▬▬"

	tmp := []rune(str)
	tmp[int(emoPos)] = '🔘'
	str = string(tmp)
	return strings.Join([]string{str, fmt.Sprintf("%s / %s", nowTime.String(), totalTime.String())}, " ")
}
func init() {
	Commands = append(Commands, &bot.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        "nowplaying",
			Description: "WHAT IS NOW PLAYING IS SHIIIIITTTTTTT (only applies to underworld)",
		},
		Run: func(b *bot.Bot, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
			vc, ok := b.VoiceInstances[evt.GuildID]
			if !ok || vc.NowPlaying == nil || vc.NowPlaying.Track == nil { // not sure if this could happen but just in case
				return &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "Error",
							Description: "not doing anything rn...",
							Color:       0xff0000,
						},
					},
				}
			}

			embed := &discordgo.MessageEmbed{
				Title:       vc.NowPlaying.Track.Title,
				Description: DescProgress(uint(vc.PlaybackPosition), vc.NowPlaying.Track.Length),
				Author: &discordgo.MessageEmbedAuthor{
					Name:    vc.NowPlaying.Member.User.Username,
					IconURL: vc.NowPlaying.Member.AvatarURL("1024"),
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text: vc.NowPlaying.Track.Author,
				},
			}
			return &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					embed,
				},
			}
		},
	})
}
