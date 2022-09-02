package bot

import (
	"math/rand"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/lukasl-dev/waterlink/v2/track"
)

type Queue struct {
	Member *discordgo.Member
	Tracks []track.Track
	mu     sync.Mutex
}

func (q *Queue) Add(track track.Track) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.Tracks = append(q.Tracks, track)
}

func (q *Queue) AddFront(song track.Track) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.Tracks = append([]track.Track{song}, q.Tracks...)
}

func (q *Queue) Pop() *track.Track {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.Tracks) == 0 {
		return nil
	}
	track := q.Tracks[0]
	q.Tracks = q.Tracks[1:]
	return &track
}

func (q *Queue) Shuffle() {
	q.mu.Lock()
	defer q.mu.Unlock()
	for i := range q.Tracks {
		j := rand.Intn(i + 1)
		q.Tracks[i], q.Tracks[j] = q.Tracks[j], q.Tracks[i]
	}
}
