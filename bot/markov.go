package bot

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/mb-14/gomarkov"
)

type Markov struct {
	chain *gomarkov.Chain
	file  *os.File
	lock  *sync.Mutex // i think there is issue with way library does thibngs????
}

func NewMarkov(gID string) *Markov {
	markov := &Markov{
		chain: gomarkov.NewChain(2),
		lock:  &sync.Mutex{},
	}

	file, err := os.OpenFile(fmt.Sprintf("models/%s", gID), os.O_RDWR|os.O_APPEND, 0666)
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
