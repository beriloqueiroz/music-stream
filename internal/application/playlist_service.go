package application

import (
	"context"
	"errors"
	"time"

	"github.com/beriloqueiroz/music-stream/internal/helper"
	"github.com/beriloqueiroz/music-stream/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlaylistService struct {
	db            *mongo.Database
	playlistsColl *mongo.Collection
}

func NewPlaylistService(db *mongo.Database) *PlaylistService {
	return &PlaylistService{
		db:            db,
		playlistsColl: db.Collection("playlists"),
	}
}

// make a crud
func (s *PlaylistService) CreatePlaylist(ctx context.Context, name string, ownerID string) (*models.Playlist, error) {
	if name == "" || ownerID == "" {
		return nil, errors.New("name and ownerID are required")
	}
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
	playlist.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return playlist, nil
}

func (s *PlaylistService) GetPlaylist(ctx context.Context, id string, ownerID string) (*models.Playlist, error) {
	if id == "" || ownerID == "" {
		return nil, errors.New("id and ownerID are required")
	}
	primitiveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	playlist := &models.Playlist{}
	err = s.playlistsColl.FindOne(ctx, bson.M{"_id": primitiveID, "owner_id": ownerID}).Decode(playlist)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

func (s *PlaylistService) UpdatePlaylist(ctx context.Context, id string, name string, ownerID string) (*models.Playlist, error) {
	if id == "" || ownerID == "" {
		return nil, errors.New("id and ownerID are required")
	}
	playlist := &models.Playlist{}
	primitiveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	result := s.playlistsColl.FindOneAndUpdate(ctx, bson.M{"_id": primitiveID, "owner_id": ownerID}, bson.M{"$set": bson.M{"name": name}})
	if result.Err() != nil {
		return nil, result.Err()
	}
	return playlist, nil
}

func (s *PlaylistService) DeletePlaylist(ctx context.Context, id string, ownerID string) error {
	if id == "" || ownerID == "" {
		return errors.New("id and ownerID are required")
	}
	primitiveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	result := s.playlistsColl.FindOneAndDelete(ctx, bson.M{"_id": primitiveID, "owner_id": ownerID})
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (s *PlaylistService) AddMusicToPlaylist(ctx context.Context, playlistID string, musicID string, ownerID string) error {
	if playlistID == "" || musicID == "" || ownerID == "" {
		return errors.New("playlistID, musicID and ownerID are required")
	}
	playlist := &models.Playlist{}
	primitiveID, err := primitive.ObjectIDFromHex(playlistID)
	if err != nil {
		return err
	}
	err = s.playlistsColl.FindOne(ctx, bson.M{"_id": primitiveID, "owner_id": ownerID}).Decode(playlist)
	if err != nil {
		return err
	}
	// verify if music already exists in playlist
	for _, music := range playlist.Musics {
		if music.MusicID == musicID {
			return errors.New("music already exists in playlist")
		}
	}
	playlist.Musics = append(playlist.Musics, models.PlaylistMusic{
		ID:         primitive.NewObjectID().Hex(),
		PlaylistID: playlistID,
		MusicID:    musicID,
		CreatedAt:  time.Now(),
	})
	_, err = s.playlistsColl.UpdateOne(ctx, bson.M{"_id": primitiveID, "owner_id": ownerID}, bson.M{"$set": bson.M{"musics": playlist.Musics}})
	if err != nil {
		return err
	}
	return nil
}

func (s *PlaylistService) RemoveMusicFromPlaylist(ctx context.Context, playlistID string, musicID string, ownerID string) error {
	if playlistID == "" || musicID == "" || ownerID == "" {
		return errors.New("playlistID, musicID and ownerID are required")
	}
	primitiveID, err := primitive.ObjectIDFromHex(playlistID)
	if err != nil {
		return err
	}
	playlist := &models.Playlist{}
	err = s.playlistsColl.FindOne(ctx, bson.M{"_id": primitiveID, "owner_id": ownerID}).Decode(playlist)
	if err != nil {
		return err
	}
	playlist.Musics = helper.RemoveFromSlice(playlist.Musics, func(music models.PlaylistMusic) bool {
		return music.MusicID == musicID
	})
	_, err = s.playlistsColl.UpdateOne(ctx, bson.M{"_id": primitiveID, "owner_id": ownerID}, bson.M{"$set": bson.M{"musics": playlist.Musics}})
	if err != nil {
		return err
	}
	return nil
}

func (s *PlaylistService) GetPlaylists(ctx context.Context, ownerID string) ([]*models.Playlist, error) {
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
