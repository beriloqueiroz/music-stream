package mongodb

import (
	"time"

	"github.com/beriloqueiroz/music-stream/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// playlist
type MongoPlaylist struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name      string               `bson:"name" json:"name"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
	Musics    []MongoPlaylistMusic `bson:"musics" json:"musics"`
	OwnerID   string               `bson:"owner_id" json:"owner_id"`
}

type MongoPlaylistMusic struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PlaylistID primitive.ObjectID `bson:"playlist_id" json:"playlist_id"`
	MusicID    primitive.ObjectID `bson:"music_id" json:"music_id"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

func (p *MongoPlaylist) ToModel() *models.Playlist {
	return &models.Playlist{
		ID:        p.ID.Hex(),
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (p *MongoPlaylistMusic) ToModel() *models.PlaylistMusic {
	return &models.PlaylistMusic{
		ID:         p.ID.Hex(),
		PlaylistID: p.PlaylistID.Hex(),
		MusicID:    p.MusicID.Hex(),
	}
}

// by model

func (p *MongoPlaylist) ByModel(model *models.Playlist) {
	id, err := primitive.ObjectIDFromHex(model.ID)
	if err != nil {
		return
	}
	p.ID = id
	p.CreatedAt = model.CreatedAt
	p.UpdatedAt = model.UpdatedAt
	p.Name = model.Name
	p.OwnerID = model.OwnerID
}

func (p *MongoPlaylistMusic) ByModel(model *models.PlaylistMusic) {
	id, err := primitive.ObjectIDFromHex(model.ID)
	if err != nil {
		return
	}
	p.ID = id
	p.CreatedAt = model.CreatedAt
	playlistID, err := primitive.ObjectIDFromHex(model.PlaylistID)
	if err != nil {
		return
	}
	p.PlaylistID = playlistID
	musicID, err := primitive.ObjectIDFromHex(model.MusicID)
	if err != nil {
		return
	}
	p.MusicID = musicID
}
