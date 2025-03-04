package application

import (
	"context"
	"errors"
	"log"
	"time"

	domain "github.com/beriloqueiroz/music-stream/internal/domain/entities"
	"github.com/beriloqueiroz/music-stream/internal/helper"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	primitiveOwnerID, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return nil, err
	}

	playlist := &domain.Playlist{
		ID:        primitive.NewObjectID(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Musics:    []domain.PlaylistMusic{},
		OwnerID:   primitiveOwnerID,
	}
	id, err := s.playlistRepo.Create(ctx, playlist)

	if err != nil {
		return nil, err
	}

	if playlist.ID.Hex() != id {
		log.Println("failed to create playlist")
		return nil, errors.New("failed to create playlist")
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
		if music.MusicID.Hex() == musicID {
			return errors.New("music already exists in playlist")
		}
	}

	// todo find with musicRepo the music by id
	primitiveMusicID, err := primitive.ObjectIDFromHex(musicID)
	if err != nil {
		return err
	}

	playlist.Musics = append(playlist.Musics, domain.PlaylistMusic{
		PlaylistID: playlist.ID,
		MusicID:    primitiveMusicID,
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
		return music.MusicID.Hex() == musicID
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
