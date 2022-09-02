package bot

import (
	"errors"

	"go.etcd.io/bbolt"
)

func (b *Bot) AddPlaylist(playlistName string, playlistURL string) error {
	return b.DB.Update(func(tx *bbolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte("playlists"))
		if err != nil {
			return err
		}
		bk.Put([]byte(playlistName), []byte(playlistURL))
		return nil
	})
}

func (b *Bot) GetPlaylist(playlistName string) (string, error) {
	var playlistURL string
	err := b.DB.View(func(tx *bbolt.Tx) error {
		bk := tx.Bucket([]byte("playlists"))
		if bk == nil {
			return errors.New("playlist bucket no made")
		}

		val := bk.Get([]byte(playlistName))
		if len(val) == 0 {
			return errors.New("playlist not found")
		}
		playlistURL = string(val)
		return nil
	})
	return playlistURL, err
}

func (b *Bot) RemovePlaylist(playlistName string) error {
	return b.DB.Update(func(tx *bbolt.Tx) error {
		bk := tx.Bucket([]byte("playlists"))
		if bk == nil {
			return errors.New("playlist bucket no made")
		}

		err := bk.Delete([]byte(playlistName))

		return err
	})
}

func (b *Bot) GetAllPlaylists() (map[string]string, error) {
	var playlists = make(map[string]string)
	err := b.DB.View(func(tx *bbolt.Tx) error {
		bk := tx.Bucket([]byte("playlists"))
		if bk == nil {
			return errors.New("playlist bucket no made")
		}
		bk.ForEach(func(k, v []byte) error {
			playlists[string(k)] = string(v)
			return nil
		})
		return nil
	})
	return playlists, err
}
