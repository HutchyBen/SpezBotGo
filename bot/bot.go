package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/gompus/snowflake"
	"github.com/lukasl-dev/waterlink/v2"
	"github.com/lukasl-dev/waterlink/v2/event"
)

type Bot struct {
	Config         *Config
	Client         *discordgo.Session
	LLConn         *waterlink.Connection
	LLClient       *waterlink.Client
	CH             *CommandHandler
	VoiceInstances map[string]*VoiceInstance
}

func NewBot(configPath string) (*Bot, error) {
	var bot Bot
	var err error
	bot.Config, err = NewConfig(configPath)
	if err != nil {
		return nil, err
	}
	bot.Client, err = discordgo.New("Bot " + bot.Config.DToken)
	if err != nil {
		return nil, err
	}
	err = bot.Client.Open()
	if err != nil {
		return nil, err
	}
	creds := waterlink.Credentials{
		Authorization: bot.Config.LLPassword,
		UserID:        snowflake.MustParse(bot.Client.State.User.ID),
	}
	bot.LLClient, err = waterlink.NewClient("http://localhost:2333", creds)
	if err != nil {
		return nil, err
	}
	bot.LLConn, err = waterlink.Open("ws://localhost:2333", creds, waterlink.ConnectionOptions{
		EventHandler: waterlink.EventHandlerFunc(bot.WaterLinkEventHandler),
	})
	if err != nil {
		return nil, err
	}
	bot.CH = NewCH()

	bot.Client.AddHandler(bot.HandleInteraction)
	bot.Client.AddHandler(bot.VoiceServerUpdate)
	bot.VoiceInstances = make(map[string]*VoiceInstance)
	return &bot, nil
}

func (b *Bot) VoiceServerUpdate(s *discordgo.Session, evt *discordgo.VoiceServerUpdate) {
	guild := b.LLConn.Guild(snowflake.MustParse(evt.GuildID))
	err := guild.UpdateVoice(s.State.SessionID, evt.Token, evt.Endpoint)
	if err != nil {
		fmt.Println("Error updating voice server:", err)
	}
}

func (b *Bot) WaterLinkEventHandler(evt interface{}) {
	switch e := evt.(type) {
	case event.TrackEnd:
		vi, ok := b.VoiceInstances[e.GuildID.String()]
		if !ok {
			return
		}
		vi.TrackEnd(e)
		break
	case event.TrackStart:
		vi, ok := b.VoiceInstances[e.GuildID.String()]
		if !ok {
			return
		}
		vi.TrackStart(e)
		break
	}

}
