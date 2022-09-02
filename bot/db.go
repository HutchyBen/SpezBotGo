package bot

import (
	"github.com/xujiajun/nutsdb"
)

func (b *Bot) AddPlaylist(playlistName string, playlistURL string) error {
	return b.DB.Update(func(tx *nutsdb.Tx) error {
		return tx.Put("playlists", []byte(playlistName), []byte(playlistURL), 0)
	})
}

func (b *Bot) GetPlaylist(playlistName string) (string, error) {
	var playlistURL string
	err := b.DB.View(func(tx *nutsdb.Tx) error {
		val, err := tx.Get("playlists", []byte(playlistName))
		if err != nil {
			return err
		}
		playlistURL = string(val.Value)
		return nil
	})
	return playlistURL, err
}

func (b *Bot) RemovePlaylist(playlistName string) error {
	return b.DB.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete("playlists", []byte(playlistName))
	})
}

func (b *Bot) GetAllPlaylists() (map[string]string, error) {
	var playlists = make(map[string]string)
	err := b.DB.View(func(tx *nutsdb.Tx) error {
		entries, err := tx.GetAll("playlists")
		if err != nil {
			return err
		}
		for _, entry := range entries {
			playlists[string(entry.Key)] = string(entry.Value)
		}
		return nil
	})
	return playlists, err
}
