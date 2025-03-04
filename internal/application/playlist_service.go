package application

import (
	"context"
	"errors"
	"time"

	domain "github.com/beriloqueiroz/music-stream/internal/domain/entities"
	"github.com/beriloqueiroz/music-stream/internal/helper"
)

type PlaylistService struct {
	playlistRepo PlaylistRepository
}

func NewPlaylistService(playlistRepo PlaylistRepository) *PlaylistService {
	return &PlaylistService{
		playlistRepo: playlistRepo,
	}
}

// make a crud
func (s *PlaylistService) CreatePlaylist(ctx context.Context, name string, ownerID string) (*domain.Playlist, error) {
	if name == "" || ownerID == "" {
		return nil, errors.New("name and ownerID are required")
	}
	playlist := &domain.Playlist{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Musics:    []domain.PlaylistMusic{},
		OwnerID:   ownerID,
	}
	id, err := s.playlistRepo.Create(ctx, playlist)

	playlist.ID = id

	if err != nil {
		return nil, err
	}
	return playlist, nil
}

func (s *PlaylistService) GetPlaylist(ctx context.Context, id string, ownerID string) (*domain.Playlist, error) {
	if id == "" || ownerID == "" {
		return nil, errors.New("id and ownerID are required")
	}
	playlist, err := s.playlistRepo.FindByID(ctx, id, ownerID)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

func (s *PlaylistService) UpdatePlaylist(ctx context.Context, id string, name string, ownerID string) (*domain.Playlist, error) {
	if id == "" || ownerID == "" {
		return nil, errors.New("id and ownerID are required")
	}
	playlist := &domain.Playlist{}
	err := s.playlistRepo.Update(ctx, playlist)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

func (s *PlaylistService) DeletePlaylist(ctx context.Context, id string, ownerID string) error {
	if id == "" || ownerID == "" {
		return errors.New("id and ownerID are required")
	}
	err := s.playlistRepo.Delete(ctx, id, ownerID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PlaylistService) AddMusicToPlaylist(ctx context.Context, playlistID string, musicID string, ownerID string) error {
	if playlistID == "" || musicID == "" || ownerID == "" {
		return errors.New("playlistID, musicID and ownerID are required")
	}
	playlist, err := s.playlistRepo.FindByID(ctx, playlistID, ownerID)
	if err != nil {
		return err
	}
	// verify if music already exists in playlist
	for _, music := range playlist.Musics {
		if music.MusicID == musicID {
			return errors.New("music already exists in playlist")
		}
	}
	playlist.Musics = append(playlist.Musics, domain.PlaylistMusic{
		PlaylistID: playlistID,
		MusicID:    musicID,
		CreatedAt:  time.Now(),
	})
	err = s.playlistRepo.Update(ctx, playlist)
	if err != nil {
		return err
	}
	return nil
}

func (s *PlaylistService) RemoveMusicFromPlaylist(ctx context.Context, playlistID string, musicID string, ownerID string) error {
	if playlistID == "" || musicID == "" || ownerID == "" {
		return errors.New("playlistID, musicID and ownerID are required")
	}
	playlist, err := s.playlistRepo.FindByID(ctx, playlistID, ownerID)
	if err != nil {
		return err
	}
	playlist.Musics = helper.RemoveFromSlice(playlist.Musics, func(music domain.PlaylistMusic) bool {
		return music.MusicID == musicID
	})
	err = s.playlistRepo.Update(ctx, playlist)
	if err != nil {
		return err
	}
	return nil
}

func (s *PlaylistService) GetPlaylists(ctx context.Context, ownerID string, page int, limit int) ([]*domain.Playlist, error) {
	playlists, err := s.playlistRepo.List(ctx, ownerID, page, limit)
	if err != nil {
		return nil, err
	}
	return playlists, nil
}
