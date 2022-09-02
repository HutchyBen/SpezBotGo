package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type InteractionCreate = func(*Bot, *discordgo.InteractionCreate) *discordgo.InteractionResponseData

type Command struct {
	discordgo.ApplicationCommand
	Run InteractionCreate
}

type CommandHandler struct {
	Commands map[string]*Command
	Events   map[string]chan *discordgo.InteractionCreate //gID, eventName
}

func NewCH() *CommandHandler {
	return &CommandHandler{Commands: make(map[string]*Command, 0), Events: make(map[string]chan *discordgo.InteractionCreate, 0)}
}

func (c *CommandHandler) Add(cmd Command) {
	c.Commands[cmd.ApplicationCommand.Name] = &cmd
}

func (c *CommandHandler) AddBulk(cmds []*Command) {
	for _, cmd := range cmds {
		c.Add(*cmd)
	}
}

func (c *CommandHandler) Register(s *discordgo.Session) {
	commands := make([]*discordgo.ApplicationCommand, 0)
	for _, cmd := range c.Commands {
		commands = append(commands, &cmd.ApplicationCommand)
	}
	s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", commands)
}

func (b *Bot) HandleInteraction(s *discordgo.Session, evt *discordgo.InteractionCreate) {
	switch evt.Type {
	case discordgo.InteractionApplicationCommand:
		cmd := b.CH.Commands[evt.ApplicationCommandData().Name]
		if cmd == nil {

			s.InteractionRespond(evt.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "I don't know that command",
				},
			})
			return
		}

		// asume response is handled if is nil

		err := s.InteractionRespond(evt.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: cmd.Run(b, evt),
		})
		if err != nil {
			fmt.Println(err)
		}
	case discordgo.InteractionMessageComponent:
		b.CH.Events[evt.GuildID] <- evt
	}
}
