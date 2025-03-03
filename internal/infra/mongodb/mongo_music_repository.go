package mongodb

import (
	"context"

	"github.com/beriloqueiroz/music-stream/internal/application"
	"github.com/beriloqueiroz/music-stream/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoMusicRepository struct {
	db         *mongo.Database
	musicsColl *mongo.Collection
}

func NewMongoMusicRepository(db *mongo.Database) *MongoMusicRepository {
	return &MongoMusicRepository{db: db, musicsColl: db.Collection("musics")}
}

func (r *MongoMusicRepository) Create(ctx context.Context, music *models.Music) (string, error) {
	mongoMusic := &MongoMusic{}
	mongoMusic.ByModel(music)
	result, err := r.musicsColl.InsertOne(ctx, mongoMusic)
	if err != nil {
		return "", err
	}
	primitiveID := result.InsertedID.(primitive.ObjectID)
	return primitiveID.Hex(), nil
}

func (r *MongoMusicRepository) FindByID(ctx context.Context, id string) (*models.Music, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	music := &models.Music{}
	err = r.musicsColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(music)
	if err != nil {
		return nil, err
	}

	return music, nil
}

func (r *MongoMusicRepository) Search(ctx context.Context, query string, page int, limit int) (*application.SearchResult, error) {
	musics, err := r.musicsColl.Find(ctx, bson.M{
		"$text": bson.M{"$search": query},
	}, options.Find().SetSkip(int64(page*limit)).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}

	var musicsList []*models.Music
	for musics.Next(ctx) {
		var music = &models.Music{}
		if err := musics.Decode(&music); err != nil {
			return nil, err
		}
		musicsList = append(musicsList, music)
	}

	return &application.SearchResult{
		MusicList: musicsList,
		Total:     len(musicsList),
	}, nil

}
