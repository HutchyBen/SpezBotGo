package bot

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/gompus/snowflake"
	"github.com/lukasl-dev/waterlink/v2"
	"github.com/lukasl-dev/waterlink/v2/event"
	"go.etcd.io/bbolt"
)

type Bot struct {
	Config         *Config
	Client         *discordgo.Session
	LLConn         *waterlink.Connection
	LLClient       *waterlink.Client
	CH             *CommandHandler
	VoiceInstances map[string]*VoiceInstance
	DB             *bbolt.DB
	Markov         map[string]*Markov
	VoiceStati     map[string]*VoiceStatus
	Die            chan (bool)
}

type VoiceStatus struct {
	token    string
	endpoint string
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
	bot.VoiceInstances = make(map[string]*VoiceInstance)
	bot.DB, err = bbolt.Open("spez.db", 0600, nil)
	bot.VoiceStati = make(map[string]*VoiceStatus)
	bot.Markov = make(map[string]*Markov)
	bot.Die = make(chan bool)

	bot.LoadMarkovChainsFromDir("models")

	bot.InitHandlers()
	return &bot, err
}

func (b *Bot) InitHandlers() {
	b.Client.AddHandler(b.VoiceServerUpdate)
	b.Client.AddHandler(b.HandleInteraction)
	b.Client.AddHandler(b.MarkovMessage)
	b.Client.AddHandler(b.VoiceStateChange)
	b.Client.AddHandler(b.E621Censor)
}

func (b *Bot) Wait() {
	b.InitDumbStuff()
	<-b.Die
}

func (b *Bot) VoiceServerUpdate(s *discordgo.Session, evt *discordgo.VoiceServerUpdate) {
	b.VoiceStati[evt.GuildID] = &VoiceStatus{
		token:    evt.Token,
		endpoint: evt.Endpoint,
	}

	b.LLConn.Guild(snowflake.MustParse(evt.GuildID)).UpdateVoice(b.Client.State.SessionID, evt.Token, evt.Endpoint)

}

func (b *Bot) WaterLinkEventHandler(evt interface{}) {
	switch e := evt.(type) {
	case event.TrackEnd:
		vi, ok := b.VoiceInstances[e.GuildID.String()]
		if !ok {
			return
		}
		vi.TrackEnd(e)

	case event.TrackStart:
		vi, ok := b.VoiceInstances[e.GuildID.String()]
		if !ok {
			return
		}
		vi.TrackStart(e)

	case event.PlayerUpdate:
		vi, ok := b.VoiceInstances[e.GuildID.String()]
		if !ok {
			return
		}
		vi.PlaybackPosition = e.State.Position
	}

}

func (b *Bot) GetVoiceState(guildID string) (*discordgo.VoiceState, error) {
	guild, err := b.Client.State.Guild(guildID)
	if err != nil {
		return nil, err
	}
	for _, vs := range guild.VoiceStates {
		if vs.UserID == b.Client.State.User.ID {
			return vs, nil
		}
	}

	return nil, errors.New("no voice state found")
}

func (b *Bot) VoiceStateChange(s *discordgo.Session, evt *discordgo.VoiceStateUpdate) {
	vi, ok := b.VoiceInstances[evt.GuildID]
	if !ok {
		return
	}
	if evt.UserID == s.State.User.ID && evt.ChannelID == "" {
		// bye...
		vi.Suicide()
	}
	// if spez is the only exister. Pause.
	g, err := s.State.Guild(evt.GuildID)
	if err != nil {
		return
	}
	spezVC, err := b.GetVoiceState(evt.GuildID)
	if err != nil {
		return
	}
	for _, vs := range g.VoiceStates {
		if vs.UserID != s.State.User.ID && spezVC.ChannelID == vs.ChannelID {
			vi.Guild.SetPaused(false)
			return
		}
	}

	vi.Guild.SetPaused(true)
	s.ChannelMessageSendEmbed(vi.MsgChannel, &discordgo.MessageEmbed{
		Title: "Paused due to everyone fucking off",
		Color: 0xffff00,
	})

}
