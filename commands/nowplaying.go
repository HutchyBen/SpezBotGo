// package commands

// import (
// 	"spezbot/bot"
// 	"math"
// 	"github.com/bwmarrin/discordgo"
// )

// func DescProgress()

// func init() {
// 	Commands = append(Commands, &bot.Command{
// 		ApplicationCommand: discordgo.ApplicationCommand{
// 			Type:        discordgo.ChatApplicationCommand,
// 			Name:        "nowplaying",
// 			Description: "WHAT IS NOW PLAYING IS SHIIIIITTTTTTT (only applies to underworld)",
// 		},

// 		Run: func(b *bot.Bot, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
// 			vc, ok := b.VoiceInstances[evt.GuildID]
// 			if !ok {
// 				return &discordgo.InteractionResponseData{
// 					Embeds: []*discordgo.MessageEmbed{
// 						{
// 							Title:       "Error",
// 							Description: "not doing anything rn...",
// 							Color:       0xff0000,
// 						},
// 					},
// 				}
// 			}

// 			var prog = float64(vc.PlaybackPosition) / float64(vc.NowPlaying.Track.Length);
//             var emoPos = math.Round(11 * prog);

//             var seconds = TimeSpan.FromSeconds(Math.Round(connection.CurrentState.PlaybackPosition.TotalSeconds));
//             var str = new StringBuilder("▬▬▬▬▬▬▬▬▬▬▬");
//             str.Remove(emoPos, 1);
//             str.Insert(emoPos, "🔘");
//             str.Append($"`{seconds:hh\\:mm\\:ss}/{song.Length:hh\\:mm\\:ss}`");

//             return str.ToString();

// 			embed := &discordgo.MessageEmbed{
// 				Title:       vc.NowPlaying.Track.Title,
// 				Description: ,
// 				Author: &discordgo.MessageEmbedAuthor{
// 					Name:    vc.NowPlaying.Member.User.Username,
// 					IconURL: vc.NowPlaying.Member.AvatarURL(1024),
// 				},
// 				Footer: &discordgo.MessageEmbedFooter{
// 					Text: vc.NowPlaying.Track.Author,
// 				},

// 			return &discordgo.InteractionResponseData{
// 				Embeds: []*discordgo.MessageEmbed{
// 					{
// 						Title:       "Success",
// 						Description: "I'm dead wtf....",
// 						Color:       0x00ff00,
// 					},
// 				},
// 			}
// 		},
// 	})
// }
