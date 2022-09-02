package commands

import (
	"spezbot/bot"

	"github.com/bwmarrin/discordgo"
)

func init() {
	Commands = append(Commands, &bot.Command{
		ApplicationCommand: discordgo.ApplicationCommand{
			Name:        "playlist",
			Description: "Manage playlists",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "add",
					Description: "Add a playlist",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Name of the playlist",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "url",
							Description: "URL of the playlist",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "remove",
					Description: "Remove a playlist",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "Name of the playlist",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "list",
					Description: "List all playlists",
				},
			},
		},
		Run: func(b *bot.Bot, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
			switch evt.ApplicationCommandData().Options[0].Name {
			case "add":
				return Add(b, evt)
			case "remove":
				return Remove(b, evt)
			case "list":
				return List(b, evt)
			}

			return nil
		},
	})
}

func Add(b *bot.Bot, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
	err := b.AddPlaylist(evt.ApplicationCommandData().Options[0].Options[0].StringValue(), evt.ApplicationCommandData().Options[0].Options[1].StringValue())
	if err != nil {
		return &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Error adding playlist",
					Color: 0xff0000,
				},
			},
		}
	}

	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Playlist added",
				Color: 0x00ff00,
			},
		},
	}

}

func Remove(b *bot.Bot, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
	err := b.RemovePlaylist(evt.ApplicationCommandData().Options[0].Options[0].StringValue())
	if err != nil {
		return &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Error removing playlist",
					Color: 0xff0000,
				},
			},
		}
	}

	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Playlist removed",
				Color: 0x00ff00,
			},
		},
	}
}

func List(b *bot.Bot, evt *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
	pls, err := b.GetAllPlaylists()
	if err != nil {
		return &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Error getting playlists",
					Color: 0xff0000,
				},
			},
		}
	}
	embed := &discordgo.MessageEmbed{
		Title:  "Playlists",
		Color:  0x00ff00,
		Fields: []*discordgo.MessageEmbedField{},
	}
	for name, url := range pls {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:  name,
			Value: url,
		})
	}
	return &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{
			embed,
		},
	}
}
