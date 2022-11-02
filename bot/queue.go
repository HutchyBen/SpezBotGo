package bot

import (
	"math/rand"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/lukasl-dev/waterlink/v2/track"
)

type UserTrack struct {
	Member *discordgo.Member
	*track.Track
}

type Queue struct {
	Member *discordgo.Member
	Tracks []UserTrack
	mu     *sync.Mutex
}

func (q *Queue) Add(track track.Track, m *discordgo.Member) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.Tracks = append(q.Tracks, UserTrack{
		Member: m,
		Track:  &track,
	})
}

func (q *Queue) AddFront(track track.Track, m *discordgo.Member) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.Tracks = append([]UserTrack{{
		Member: m,
		Track:  &track,
	}}, q.Tracks...)
}

func (q *Queue) Pop() *UserTrack {
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
