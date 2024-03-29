package bot

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/mb-14/gomarkov"
)

type Markov struct {
	chain *gomarkov.Chain
	file  *os.File
	lock  *sync.Mutex // i think there is issue with way library does thibngs????
}

func NewMarkov(gID string) *Markov {
	markov := &Markov{
		chain: gomarkov.NewChain(1),
		lock:  &sync.Mutex{},
	}

	file, err := os.OpenFile(fmt.Sprintf("models/%s", gID), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	fmt.Println(err)
	if err != nil {
		return markov
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return markov
	}

	for _, line := range strings.Split(string(data), "\n") {
		markov.chain.Add(strings.Split(line, " "))
	}
	markov.file = file
	return markov
}

func (m *Markov) Add(text string) {
	m.lock.Lock()
	m.chain.Add(strings.Split(text, " "))
	m.lock.Unlock()
	m.file.WriteString(text + "\n")
}

func (m *Markov) Generate() string {
	tokens := []string{gomarkov.StartToken}
	for tokens[len(tokens)-1] != gomarkov.EndToken {
		m.lock.Lock()
		next, _ := m.chain.Generate(tokens[(len(tokens) - 1):])
		m.lock.Unlock()
		tokens = append(tokens, next)
	}
	return strings.Join(tokens[1:len(tokens)-1], " ")
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

	mk, ok := b.Markov[evt.GuildID]
	if !ok {
		b.Markov[evt.GuildID] = NewMarkov(evt.GuildID)
		mk = b.Markov[evt.GuildID]
	}
	if rand.Intn(8) == 1 || strings.Contains(strings.ToLower(evt.Content), "spez") || strings.Contains(evt.Content, fmt.Sprintf("<@%s>", s.State.User.ID)) {
		s.ChannelMessageSend(evt.ChannelID, mk.Generate())
	}
	if strings.TrimSpace(evt.Content) != "" && evt.Author.ID != "947332449854193696" && strings.ToLower(evt.Content) != "spez" {
		mk.Add(strings.TrimSpace(evt.Content))
	}

}
