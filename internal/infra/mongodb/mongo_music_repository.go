package mongodb

import (
	"context"

	"github.com/beriloqueiroz/music-stream/internal/application"
	domain "github.com/beriloqueiroz/music-stream/internal/domain/entities"
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

func (r *MongoMusicRepository) Create(ctx context.Context, music *domain.Music) (string, error) {
	result, err := r.musicsColl.InsertOne(ctx, music)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *MongoMusicRepository) FindByID(ctx context.Context, id string) (*domain.Music, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	music := &domain.Music{}
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

	var musicsList []*domain.Music
	for musics.Next(ctx) {
		var music = &domain.Music{}
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
