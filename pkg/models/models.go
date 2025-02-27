package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Music struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string            `bson:"title" json:"title"`
	Artist    string            `bson:"artist" json:"artist"`
	Album     string            `bson:"album" json:"album"`
	Duration  int32             `bson:"duration" json:"duration"`
	FilePath  string            `bson:"file_path" json:"file_path"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string            `bson:"email" json:"email"`
	Password  string            `bson:"password" json:"-"`
	IsAdmin   bool              `bson:"is_admin" json:"is_admin"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
}

type Invite struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code      string            `bson:"code" json:"code"`
	Email     string            `bson:"email" json:"email"`
	Used      bool              `bson:"used" json:"used"`
	ExpiresAt time.Time         `bson:"expires_at" json:"expires_at"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
} 