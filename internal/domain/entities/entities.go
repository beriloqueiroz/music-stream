package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Music struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Artist    string             `bson:"artist" json:"artist"`
	Album     string             `bson:"album" json:"album"`
	StorageID primitive.ObjectID `bson:"storage_id" json:"storage_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Lyrics    *Lyrics            `bson:"lyrics,omitempty" json:"lyrics,omitempty"`
	Tablature *Tablature         `bson:"tablature,omitempty" json:"tablature,omitempty"`
	Type      string             `bson:"type" json:"type"`
	Duration  float64            `bson:"duration" json:"duration"`
	Metadata  *MusicMetadata     `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

type MusicMetadata struct {
	Title    string `bson:"title,omitempty" json:"title,omitempty"`
	Artist   string `bson:"artist,omitempty" json:"artist,omitempty"`
	Album    string `bson:"album,omitempty" json:"album,omitempty"`
	Type     string `bson:"type,omitempty" json:"type,omitempty"`
	Year     int32  `bson:"year,omitempty" json:"year,omitempty"`           // Ano de lançamento
	Genre    string `bson:"genre,omitempty" json:"genre,omitempty"`         // Gênero musical
	Composer string `bson:"composer,omitempty" json:"composer,omitempty"`   // Compositor
	Label    string `bson:"label,omitempty" json:"label,omitempty"`         // Gravadora
	AlbumArt string `bson:"album_art,omitempty" json:"album_art,omitempty"` // Capa do álbum
	Comments string `bson:"comments,omitempty" json:"comments,omitempty"`   // Comentários
	Isrc     string `bson:"isrc,omitempty" json:"isrc,omitempty"`           // Código ISRC
	Url      string `bson:"url,omitempty" json:"url,omitempty"`             // URL para informações adicionais
}

type Lyrics struct {
	Text     string    `bson:"text" json:"text"`         // Letra completa
	Timing   []Segment `bson:"timing" json:"timing"`     // Temporização
	Language string    `bson:"language" json:"language"` // Idioma da letra
}

type Tablature struct {
	Content string    `bson:"content" json:"content"` // Cifra/Tablatura
	Timing  []Segment `bson:"timing" json:"timing"`   // Temporização
	Format  string    `bson:"format" json:"format"`   // Formato (chord pro, etc)
}

type Segment struct {
	Start   float64 `bson:"start" json:"start"`     // Tempo em segundos
	End     float64 `bson:"end" json:"end"`         // Tempo em segundos
	Content string  `bson:"content" json:"content"` // Trecho da letra/cifra
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	IsAdmin   bool               `bson:"is_admin" json:"is_admin"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Invite struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code      string             `bson:"code" json:"code"`
	Email     string             `bson:"email" json:"email"`
	Used      bool               `bson:"used" json:"used"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expires_at"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

// playlist
type Playlist struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Musics    []PlaylistMusic    `bson:"musics" json:"musics"`
	OwnerID   primitive.ObjectID `bson:"owner_id" json:"owner_id"`
}

type PlaylistMusic struct {
	PlaylistID primitive.ObjectID `bson:"playlist_id" json:"playlist_id"`
	MusicID    primitive.ObjectID `bson:"music_id" json:"music_id"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	Title      string             `bson:"title" json:"title"`
	Artist     string             `bson:"artist" json:"artist"`
	Album      string             `bson:"album" json:"album"`
	Duration   float64            `bson:"duration" json:"duration"`
	Type       string             `bson:"type" json:"type"`
}
