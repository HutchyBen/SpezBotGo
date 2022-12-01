package bot

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

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
	// Load discord client
	bot.Client, err = discordgo.New("Bot " + bot.Config.DToken)
	if err != nil {
		return nil, err
	}
	err = bot.Client.Open()
	if err != nil {
		return nil, err
	}
	bot.Client.AddHandler(bot.VoiceStateChange)
	// Load waterlink stuff
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
	bot.Client.AddHandler(bot.VoiceServerUpdate)

	bot.CH = NewCH()
	bot.Client.AddHandler(bot.HandleInteraction)

	bot.VoiceInstances = make(map[string]*VoiceInstance)

	bot.DB, err = bbolt.Open("spez.db", 0600, nil)

	bot.Markov = make(map[string]*Markov)
	bot.VoiceStati = make(map[string]*VoiceStatus)
	bot.LoadMarkovChainsFromDir("models")
	bot.Client.AddHandler(bot.MarkovMessage)
	return &bot, err
}

func (b *Bot) LoadMarkovChainsFromDir(dir string) error {
	fs, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range fs {
		mk := NewMarkov(f.Name())
		b.Markov[strings.Split(f.Name(), ".")[0]] = mk
	}
	return nil
}

func (b *Bot) MarkovMessage(s *discordgo.Session, evt *discordgo.MessageCreate) {
	if evt.Author.ID == s.State.User.ID {
		return
	}
	text := strings.ToLower(evt.Content)
	// lol
	if strings.Contains(text, "donkey-kong") || strings.Contains(text, "KuvDsT4sRzU") || strings.Contains(text, "donkeykong") || strings.Contains(text, "donkey kong") {
		err := s.ChannelMessageDelete(evt.ChannelID, evt.ID)
		if err != nil {
			fmt.Println("Cannot delete message: " + err.Error())
		}
	}
	mk, ok := b.Markov[evt.GuildID]
	if !ok {
		b.Markov[evt.GuildID] = NewMarkov(evt.GuildID)
		mk = b.Markov[evt.GuildID]
	}
	if strings.TrimSpace(evt.Content) != "" {
		mk.Add(strings.TrimSpace(evt.Content))
	}

	if rand.Intn(8) == 1 || strings.Contains(evt.Content, "spez") || strings.Contains(evt.Content, fmt.Sprintf("<@%s>", s.State.User.ID)) {
		s.ChannelMessageSend(evt.ChannelID, mk.Generate())
	}
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

	}
}

func (b *Bot) VoiceStateChange(s *discordgo.Session, evt *discordgo.VoiceStateUpdate) {
	if evt.UserID == s.State.User.ID && evt.ChannelID == "" {
		// bye...
		vi, ok := b.VoiceInstances[evt.GuildID]
		if !ok {
			return
		}

		vi.Suicide()

	}

}
