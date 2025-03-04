package mongodb

import (
	"time"

	domain "github.com/beriloqueiroz/music-stream/internal/domain/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoUser struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	IsAdmin   bool               `bson:"is_admin" json:"is_admin"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type MongoInvite struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code      string             `bson:"code" json:"code"`
	Email     string             `bson:"email" json:"email"`
	Used      bool               `bson:"used" json:"used"`
	ExpiresAt time.Time          `bson:"expires_at" json:"expires_at"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

func (i *MongoInvite) ToModel() *domain.Invite {
	return &domain.Invite{
		ID:        i.ID.Hex(),
		Code:      i.Code,
		Email:     i.Email,
		Used:      i.Used,
		ExpiresAt: i.ExpiresAt,
		CreatedAt: i.CreatedAt,
	}
}

func (u *MongoUser) ToModel() *domain.User {
	return &domain.User{
		ID:        u.ID.Hex(),
		Email:     u.Email,
		Password:  u.Password,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
	}
}

// by model

func (i *MongoInvite) ByModel(model *domain.Invite) {
	if model.ID != "" {
		id, err := primitive.ObjectIDFromHex(model.ID)
		if err != nil {
			return
		}
		i.ID = id
	} else {
		i.ID = primitive.NewObjectID()
	}
	i.CreatedAt = model.CreatedAt
	i.Email = model.Email
	i.Code = model.Code
	i.Used = model.Used
	i.ExpiresAt = model.ExpiresAt
}

func (u *MongoUser) ByModel(model *domain.User) {
	id, err := primitive.ObjectIDFromHex(model.ID)
	if err != nil {
		return
	}
	u.ID = id
	u.CreatedAt = model.CreatedAt
	u.Email = model.Email
	u.Password = model.Password
	u.IsAdmin = model.IsAdmin
}
