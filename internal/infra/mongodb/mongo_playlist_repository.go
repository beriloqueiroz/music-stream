package mongodb

import (
	"context"

	"github.com/beriloqueiroz/music-stream/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoPlaylistRepository struct {
	db            *mongo.Database
	playlistsColl *mongo.Collection
}

func NewMongoPlaylistRepository(db *mongo.Database) *MongoPlaylistRepository {
	return &MongoPlaylistRepository{db: db, playlistsColl: db.Collection("playlists")}
}

func (r *MongoPlaylistRepository) Create(ctx context.Context, playlist *models.Playlist) (string, error) {
	mongoPlaylist := &MongoPlaylist{}
	mongoPlaylist.ByModel(playlist)
	result, err := r.playlistsColl.InsertOne(ctx, mongoPlaylist)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *MongoPlaylistRepository) FindByID(ctx context.Context, id string) (*models.Playlist, error) {
	playlist := &MongoPlaylist{}
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.playlistsColl.FindOne(ctx, bson.M{"_id": idPrimitive}).Decode(playlist)
	if err != nil {
		return nil, err
	}
	return playlist.ToModel(), nil
}

func (r *MongoPlaylistRepository) FindByUserID(ctx context.Context, userID string) ([]*models.Playlist, error) {
	playlists := []*MongoPlaylist{}
	cursor, err := r.playlistsColl.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		playlist := &MongoPlaylist{}
		err = cursor.Decode(playlist)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}
	playlistsModels := make([]*models.Playlist, len(playlists))
	for i, playlist := range playlists {
		playlistsModels[i] = playlist.ToModel()
	}
	return playlistsModels, nil
}

func (r *MongoPlaylistRepository) Update(ctx context.Context, playlist *models.Playlist) error {
	mongoPlaylist := &MongoPlaylist{}
	mongoPlaylist.ByModel(playlist)
	_, err := r.playlistsColl.UpdateOne(ctx, bson.M{"_id": playlist.ID}, bson.M{"$set": mongoPlaylist})
	return err
}

func (r *MongoPlaylistRepository) Delete(ctx context.Context, id string) error {
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.playlistsColl.DeleteOne(ctx, bson.M{"_id": idPrimitive})
	return err
}
