package mongodb

import (
	"context"

	domain "github.com/beriloqueiroz/music-stream/internal/domain/entities"
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

func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	user := &domain.User{}
	primitiveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.usersColl.FindOne(ctx, bson.M{"_id": primitiveID}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *MongoUserRepository) Insert(ctx context.Context, user *domain.User) error {
	_, err := r.usersColl.InsertOne(ctx, user)
	return err
}

func (r *MongoUserRepository) InsertInvite(ctx context.Context, invite *domain.Invite) error {
	_, err := r.invitesColl.InsertOne(ctx, invite)
	return err
}

func (r *MongoUserRepository) UpdateInvite(ctx context.Context, invite *domain.Invite) error {
	_, err := r.invitesColl.UpdateOne(ctx, bson.M{"_id": invite.ID}, bson.M{"$set": invite})
	return err
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.usersColl.FindOne(ctx, bson.M{"email": email}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *MongoUserRepository) FindInviteByCode(ctx context.Context, code string) (*domain.Invite, error) {
	invite := &domain.Invite{}
	err := r.invitesColl.FindOne(ctx, bson.M{"code": code}).Decode(invite)
	if err != nil {
		return nil, err
	}
	return invite, nil
}
