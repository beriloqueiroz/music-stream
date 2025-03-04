package mongodb

import (
	"context"

	domain "github.com/beriloqueiroz/music-stream/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoPlaylistRepository struct {
	db            *mongo.Database
	playlistsColl *mongo.Collection
}

func NewMongoPlaylistRepository(db *mongo.Database) *MongoPlaylistRepository {
	return &MongoPlaylistRepository{db: db, playlistsColl: db.Collection("playlists")}
}

func (r *MongoPlaylistRepository) Create(ctx context.Context, playlist *domain.Playlist) (string, error) {
	result, err := r.playlistsColl.InsertOne(ctx, playlist)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *MongoPlaylistRepository) FindByID(ctx context.Context, id string, ownerID string) (*domain.Playlist, error) {
	playlist := &domain.Playlist{}
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	ownerIDPrimitive, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return nil, err
	}
	err = r.playlistsColl.FindOne(ctx, bson.M{"_id": idPrimitive, "owner_id": ownerIDPrimitive}).Decode(playlist)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

func (r *MongoPlaylistRepository) List(ctx context.Context, ownerID string, page int, limit int) ([]*domain.Playlist, error) {
	playlists := []*domain.Playlist{}
	ownerIDPrimitive, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.playlistsColl.Find(ctx, bson.M{"owner_id": ownerIDPrimitive}, options.Find().SetSkip(int64(page*limit)).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		playlist := &domain.Playlist{}
		err = cursor.Decode(playlist)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}

func (r *MongoPlaylistRepository) Update(ctx context.Context, playlist *domain.Playlist) error {
	_, err := r.playlistsColl.UpdateOne(ctx, bson.M{"_id": playlist.ID, "owner_id": playlist.OwnerID}, bson.M{"$set": playlist})
	return err
}

func (r *MongoPlaylistRepository) Delete(ctx context.Context, id string, ownerID string) error {
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	ownerIDPrimitive, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return err
	}
	_, err = r.playlistsColl.DeleteOne(ctx, bson.M{"_id": idPrimitive, "owner_id": ownerIDPrimitive})
	return err
}
