package playlist

import (
	"context"
	"time"

	"github.com/beriloqueiroz/music-stream/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	db            *mongo.Database
	playlistsColl *mongo.Collection
}

func NewPlaylistService(db *mongo.Database) *Service {
	return &Service{
		db:            db,
		playlistsColl: db.Collection("playlists"),
	}
}

// make a crud
func (s *Service) CreatePlaylist(ctx context.Context, name string, ownerID string) (*models.Playlist, error) {
	playlist := &models.Playlist{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Musics:    []models.PlaylistMusic{},
		OwnerID:   ownerID,
	}
	result, err := s.playlistsColl.InsertOne(ctx, playlist)
	if err != nil {
		return nil, err
	}
	playlist.ID = result.InsertedID.(string)
	return playlist, nil
}

func (s *Service) GetPlaylist(ctx context.Context, id string, ownerID string) (*models.Playlist, error) {
	playlist := &models.Playlist{}
	err := s.playlistsColl.FindOne(ctx, bson.M{"_id": id, "owner_id": ownerID}).Decode(playlist)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

func (s *Service) UpdatePlaylist(ctx context.Context, id string, name string, ownerID string) (*models.Playlist, error) {
	playlist := &models.Playlist{}
	result := s.playlistsColl.FindOneAndUpdate(ctx, bson.M{"_id": id, "owner_id": ownerID}, bson.M{"$set": bson.M{"name": name}})
	if result.Err() != nil {
		return nil, result.Err()
	}
	return playlist, nil
}

func (s *Service) DeletePlaylist(ctx context.Context, id string, ownerID string) error {
	result := s.playlistsColl.FindOneAndDelete(ctx, bson.M{"_id": id, "owner_id": ownerID})
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (s *Service) AddMusicToPlaylist(ctx context.Context, playlistID string, musicID string, ownerID string) error {
	_, err := s.playlistsColl.UpdateOne(ctx, bson.M{"_id": playlistID, "owner_id": ownerID}, bson.M{"$push": bson.M{"musics": bson.M{"music_id": musicID}}})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RemoveMusicFromPlaylist(ctx context.Context, playlistID string, musicID string, ownerID string) error {
	_, err := s.playlistsColl.UpdateOne(ctx, bson.M{"_id": playlistID, "owner_id": ownerID}, bson.M{"$pull": bson.M{"musics": bson.M{"music_id": musicID}}})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetPlaylists(ctx context.Context, ownerID string) ([]*models.Playlist, error) {
	playlists := []*models.Playlist{}
	cursor, err := s.playlistsColl.Find(ctx, bson.M{"owner_id": ownerID})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &playlists)
	if err != nil {
		return nil, err
	}
	return playlists, nil
}
