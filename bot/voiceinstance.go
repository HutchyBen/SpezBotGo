package bot

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/gompus/snowflake"
	"github.com/lukasl-dev/waterlink/v2"
	"github.com/lukasl-dev/waterlink/v2/event"
	"github.com/lukasl-dev/waterlink/v2/track"
	"github.com/lukasl-dev/waterlink/v2/track/query"
)

type NowPlaying struct {
	Track  *track.Info
	Member *discordgo.Member
}

type VoiceInstance struct {
	Guild      waterlink.Guild
	Queues     []Queue
	QueueIndex int
	Bot        *Bot
	NowPlaying *NowPlaying
	MsgChannel string
	GID        string
}

func (b *Bot) CreateVoiceInstance(gID string, mID string, vcID string, msgID string) error {
	var vi VoiceInstance
	vi.GID = gID

	err := b.Client.ChannelVoiceJoinManual(gID, vcID, false, false)
	if err != nil {
		panic(err)
	}
	vi.Guild = b.LLConn.Guild(snowflake.MustParse(gID))

	vi.Queues = make([]Queue, 0)
	b.VoiceInstances[gID] = &vi
	vi.MsgChannel = msgID
	vi.Bot = b
	vi.NowPlaying = nil
	return nil
}

func (vi *VoiceInstance) Suicide() {
	delete(vi.Bot.VoiceInstances, vi.GID)
}

func (vi *VoiceInstance) GetSongs(search string) (*track.LoadResult, error) {
	var result *track.LoadResult
	var err error
	if _, err = url.ParseRequestURI(search); err != nil {
		fmt.Println(search)
		query := query.YouTube(search)
		fmt.Println(query)
		result, err = vi.Bot.LLClient.LoadTracks(query)
	} else {
		result, err = vi.Bot.LLClient.LoadTracks(query.Of((search)))
	}
	return result, err
}

func (vi *VoiceInstance) QueueSong(member *discordgo.Member, song track.Track) {
	for i := 0; i < len(vi.Queues); i++ {
		if vi.Queues[i].Member.User.ID == member.User.ID {
			vi.Queues[i].Add(song)
			return
		}
	}

	vi.Queues = append(vi.Queues, Queue{
		Member: member,
		Tracks: []track.Track{song},
		mu:     &sync.Mutex{},
	})
}

func (vi *VoiceInstance) QueueSongNext(member *discordgo.Member, song track.Track) {
	for i := 0; i < len(vi.Queues); i++ {
		if vi.Queues[i].Member.User.ID == member.User.ID {
			vi.Queues[i].AddFront(song)
			vi.QueueIndex = i - 1
			return
		}
	}

	vi.Queues = append(vi.Queues, Queue{
		Member: member,
		Tracks: []track.Track{song},
	})
}

func (vi *VoiceInstance) TrackStart(evt event.TrackStart) {
	var err error

	vi.NowPlaying.Track, err = vi.Bot.LLClient.DecodeTrack(evt.TrackID)
	if err != nil {
		fmt.Println(err)
		return
	}
	vi.Bot.Client.ChannelMessageSendEmbed(vi.MsgChannel, &discordgo.MessageEmbed{
		Title:       "Now Playing",
		Description: vi.NowPlaying.Track.Title,
		Color:       0x00ff00,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    vi.NowPlaying.Member.User.Username,
			IconURL: vi.NowPlaying.Member.User.AvatarURL(""),
		},
		URL: vi.NowPlaying.Track.URI,
		Footer: &discordgo.MessageEmbedFooter{
			Text: vi.NowPlaying.Track.Author,
		},
	})
}

// Load next song
func (vi *VoiceInstance) TrackEnd(evt event.TrackEnd) {
	playbackFinishEmbed := &discordgo.MessageEmbed{
		Title: "Playback Finished",
		Color: 0xffff00,
	}

	if len(vi.Queues) == 0 {
		vi.Bot.Client.ChannelMessageSendEmbed(vi.MsgChannel, playbackFinishEmbed)
		vi.Suicide()
		return
	}

	vi.QueueIndex++
	if vi.QueueIndex >= len(vi.Queues) {
		vi.QueueIndex = 0
	}

	vi.NowPlaying.Member = vi.Queues[vi.QueueIndex].Member
	vi.Guild.PlayTrack(*vi.Queues[vi.QueueIndex].Pop())

	if len(vi.Queues[vi.QueueIndex].Tracks) == 0 {
		vi.Queues = append(vi.Queues[:vi.QueueIndex], vi.Queues[vi.QueueIndex+1:]...)
	}
}
