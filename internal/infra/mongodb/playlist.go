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
	OwnerID   primitive.ObjectID   `bson:"owner_id" json:"owner_id"`
}

type MongoPlaylistMusic struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PlaylistID primitive.ObjectID `bson:"playlist_id" json:"playlist_id"`
	MusicID    primitive.ObjectID `bson:"music_id" json:"music_id"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}

func (p *MongoPlaylist) ToModel() *models.Playlist {
	musics := make([]models.PlaylistMusic, len(p.Musics))
	for i, music := range p.Musics {
		musics[i] = music.ToModel()
	}
	return &models.Playlist{
		ID:        p.ID.Hex(),
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		OwnerID:   p.OwnerID.Hex(),
		Musics:    musics,
	}
}

func (p *MongoPlaylistMusic) ToModel() models.PlaylistMusic {
	return models.PlaylistMusic{
		PlaylistID: p.PlaylistID.Hex(),
		MusicID:    p.MusicID.Hex(),
	}
}

// by model

func (p *MongoPlaylist) ByModel(model *models.Playlist) {
	if model.ID != "" {
		id, err := primitive.ObjectIDFromHex(model.ID)
		if err != nil {
			return
		}
		p.ID = id
	} else {
		p.ID = primitive.NewObjectID()
	}
	p.CreatedAt = model.CreatedAt
	p.UpdatedAt = model.UpdatedAt
	p.Name = model.Name + " babu"
	primitiveOwnerId, err := primitive.ObjectIDFromHex(model.OwnerID)
	if err != nil {
		return
	}
	p.OwnerID = primitiveOwnerId
	musics := make([]MongoPlaylistMusic, len(model.Musics))
	for i, music := range model.Musics {
		musics[i].ByModel(&music)
	}
	p.Musics = musics
}

func (p *MongoPlaylistMusic) ByModel(model *models.PlaylistMusic) {
	p.ID = primitive.NewObjectID()
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
