package mongodb

import (
	"context"

	"github.com/beriloqueiroz/music-stream/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	db          *mongo.Database
	usersColl   *mongo.Collection
	invitesColl *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		db:          db,
		usersColl:   db.Collection("users"),
		invitesColl: db.Collection("invites"),
	}
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	user := &MongoUser{}
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.usersColl.FindOne(ctx, bson.M{"_id": idPrimitive}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user.ToModel(), nil
}

func (r *MongoUserRepository) Insert(ctx context.Context, user *models.User) error {
	mongoUser := &MongoUser{}
	mongoUser.ByModel(user)
	_, err := r.usersColl.InsertOne(ctx, mongoUser)
	return err
}

func (r *MongoUserRepository) InsertInvite(ctx context.Context, invite *models.Invite) error {
	mongoInvite := &MongoInvite{}
	mongoInvite.ByModel(invite)
	_, err := r.invitesColl.InsertOne(ctx, mongoInvite)
	return err
}

func (r *MongoUserRepository) UpdateInvite(ctx context.Context, invite *models.Invite) error {
	mongoInvite := &MongoInvite{}
	mongoInvite.ByModel(invite)
	_, err := r.invitesColl.UpdateOne(ctx, bson.M{"_id": invite.ID}, bson.M{"$set": mongoInvite})
	return err
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &MongoUser{}
	err := r.usersColl.FindOne(ctx, bson.M{"email": email}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user.ToModel(), nil
}

func (r *MongoUserRepository) FindInviteByCode(ctx context.Context, code string) (*models.Invite, error) {
	invite := &MongoInvite{}
	err := r.invitesColl.FindOne(ctx, bson.M{"code": code}).Decode(invite)
	if err != nil {
		return nil, err
	}
	return invite.ToModel(), nil
}
