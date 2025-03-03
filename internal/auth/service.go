package auth

import (
	"context"
	"errors"
	"time"

	"github.com/beriloqueiroz/music-stream/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db          *mongo.Database
	jwtSecret   []byte
	invitesColl *mongo.Collection
	usersColl   *mongo.Collection
}

func NewAuthService(db *mongo.Database, jwtSecret string) *Service {
	return &Service{
		db:          db,
		jwtSecret:   []byte(jwtSecret),
		invitesColl: db.Collection("invites"),
		usersColl:   db.Collection("users"),
	}
}

func (s *Service) CreateInvite(ctx context.Context, email string, whoIsInvitingId string) (*models.Invite, error) {
	whoIsInviting := &models.User{}
	primitiveWhoIsInvitingId, err := primitive.ObjectIDFromHex(whoIsInvitingId)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}
	err = s.usersColl.FindOne(ctx, bson.M{"_id": primitiveWhoIsInvitingId}).Decode(whoIsInviting)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}
	if !whoIsInviting.IsAdmin {
		return nil, errors.New("usuário tem permissão insuficiente para criar convite")
	}
	code := generateRandomCode() // Implementar função para gerar código aleatório
	invite := &models.Invite{
		ID:        primitive.NewObjectID().Hex(),
		Code:      code,
		Email:     email,
		Used:      false,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // Expira em 7 dias
		CreatedAt: time.Now(),
	}

	_, err = s.invitesColl.InsertOne(ctx, invite)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

func (s *Service) Register(ctx context.Context, email, password, inviteCode string) error {
	// Verificar se o convite existe e é válido
	invite := &models.Invite{}
	err := s.invitesColl.FindOne(ctx, bson.M{
		"code":       inviteCode,
		"email":      email,
		"used":       false,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(invite)

	if err != nil {
		return errors.New("convite inválido ou expirado")
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Criar usuário
	user := &models.User{
		ID:        primitive.NewObjectID().Hex(),
		Email:     email,
		Password:  string(hashedPassword),
		IsAdmin:   false,
		CreatedAt: time.Now(),
	}

	_, err = s.usersColl.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	// Marcar convite como usado
	_, err = s.invitesColl.UpdateOne(
		ctx,
		bson.M{"_id": invite.ID},
		bson.M{"$set": bson.M{"used": true}},
	)

	return err
}

func (s *Service) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	user := &models.User{}
	err := s.usersColl.FindOne(ctx, bson.M{"email": email}).Decode(user)
	if err != nil {
		return "", nil, errors.New("usuário não encontrado")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("senha incorreta")
	}

	// Gerar JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}
