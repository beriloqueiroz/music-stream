package mongodb

import (
	"time"

	"github.com/beriloqueiroz/music-stream/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoMusic struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Artist    string             `bson:"artist" json:"artist"`
	Album     string             `bson:"album" json:"album"`
	StorageID string             `bson:"storage_id" json:"storage_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Lyrics    *MongoLyrics       `bson:"lyrics,omitempty" json:"lyrics,omitempty"`
	Tablature *MongoTablature    `bson:"tablature,omitempty" json:"tablature,omitempty"`
}

type MongoLyrics struct {
	Text     string         `bson:"text" json:"text"`         // Letra completa
	Timing   []MongoSegment `bson:"timing" json:"timing"`     // Temporização
	Language string         `bson:"language" json:"language"` // Idioma da letra
}

type MongoTablature struct {
	Content string         `bson:"content" json:"content"` // Cifra/Tablatura
	Timing  []MongoSegment `bson:"timing" json:"timing"`   // Temporização
	Format  string         `bson:"format" json:"format"`   // Formato (chord pro, etc)
}

type MongoSegment struct {
	Start   float64 `bson:"start" json:"start"`     // Tempo em segundos
	End     float64 `bson:"end" json:"end"`         // Tempo em segundos
	Content string  `bson:"content" json:"content"` // Trecho da letra/cifra
}

func (m *MongoMusic) ToModel() *models.Music {
	return &models.Music{
		ID:        m.ID.Hex(),
		Title:     m.Title,
		Artist:    m.Artist,
		Album:     m.Album,
		StorageID: m.StorageID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Lyrics:    m.Lyrics.ToModel(),
		Tablature: m.Tablature.ToModel(),
	}
}

func (l *MongoLyrics) ToModel() *models.Lyrics {
	timing := make([]models.Segment, len(l.Timing))
	for i, s := range l.Timing {
		timing[i] = *s.ToModel()
	}
	return &models.Lyrics{
		Text:     l.Text,
		Timing:   timing,
		Language: l.Language,
	}
}

func (t *MongoTablature) ToModel() *models.Tablature {
	timing := make([]models.Segment, len(t.Timing))
	for i, s := range t.Timing {
		timing[i] = *s.ToModel()
	}
	return &models.Tablature{
		Content: t.Content,
		Timing:  timing,
		Format:  t.Format,
	}
}

func (s *MongoSegment) ToModel() *models.Segment {
	return &models.Segment{
		Start:   s.Start,
		End:     s.End,
		Content: s.Content,
	}
}

// by model
func (m *MongoMusic) ByModel(model *models.Music) {
	if model.ID != "" {
		id, err := primitive.ObjectIDFromHex(model.ID)
		if err != nil {
			return
		}
		m.ID = id
	} else {
		m.ID = primitive.NewObjectID()
	}
	m.CreatedAt = model.CreatedAt
	m.UpdatedAt = model.UpdatedAt
	m.Title = model.Title
	m.Artist = model.Artist
	m.Album = model.Album
	m.StorageID = model.StorageID
	if model.Lyrics != nil {
		m.Lyrics = &MongoLyrics{}
		m.Lyrics.ByModel(model.Lyrics)
	}
	if model.Tablature != nil {
		m.Tablature = &MongoTablature{}
		m.Tablature.ByModel(model.Tablature)
	}
}

func (m *MongoLyrics) ByModel(model *models.Lyrics) {
	m.Text = model.Text
	m.Timing = make([]MongoSegment, len(model.Timing))
	for i, s := range model.Timing {
		m.Timing[i] = MongoSegment{
			Start:   s.Start,
			End:     s.End,
			Content: s.Content,
		}
	}
}

func (m *MongoTablature) ByModel(model *models.Tablature) {
	m.Content = model.Content
	m.Timing = make([]MongoSegment, len(model.Timing))
	for i, s := range model.Timing {
		m.Timing[i] = MongoSegment{
			Start:   s.Start,
			End:     s.End,
			Content: s.Content,
		}
	}
}

func (m *MongoSegment) ByModel(model *models.Segment) {
	m.Start = model.Start
	m.End = model.End
	m.Content = model.Content
}
