package commands

import (
	"errors"
	"fmt"
	"net/url"
	"spezbot/bot"

	"github.com/bwmarrin/discordgo"
	"github.com/lukasl-dev/waterlink/v2/track"
)

func PlayStart(b *bot.Bot, evt *discordgo.InteractionCreate) (songs *track.LoadResult, vi *bot.VoiceInstance, isURL bool, err error) {
	vi, ok := b.VoiceInstances[evt.GuildID]
	if !ok {
		state, err := b.Client.State.VoiceState(evt.GuildID, evt.Member.User.ID)
		if err != nil {
			return nil, nil, false, errors.New("you are not in a voice chat")
		}
		fmt.Println(state.ChannelID == "")
		b.CreateVoiceInstance(evt.GuildID, evt.Member.User.ID, state.ChannelID, evt.ChannelID)
		vi = b.VoiceInstances[evt.GuildID]
	}
	songs, err = vi.GetSongs(evt.ApplicationCommandData().Options[0].StringValue())
	if err != nil {
		return nil, nil, false, errors.New("error getting songs")
	}
	if len(songs.Tracks) == 0 {
		return nil, nil, false, errors.New("no songs found")
	}
	_, uri := url.ParseRequestURI(evt.ApplicationCommandData().Options[0].StringValue())
	isURL = uri == nil
	// If no song is playing play the song
	return
}
