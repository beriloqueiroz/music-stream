package application

import (
	"context"
	"errors"
	"time"

	"github.com/beriloqueiroz/music-stream/internal/helper"
	"github.com/beriloqueiroz/music-stream/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	jwtSecret []byte
	userRepo  UserRepository
}

func NewUserService(userRepo UserRepository, jwtSecret []byte) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *UserService) CreateInvite(ctx context.Context, email string, whoIsInvitingId string) (*models.Invite, error) {
	whoIsInviting, err := s.userRepo.FindByID(ctx, whoIsInvitingId)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}
	if !whoIsInviting.IsAdmin {
		return nil, errors.New("usuário tem permissão insuficiente para criar convite")
	}
	code := helper.GenerateRandomCode() // Implementar função para gerar código aleatório
	invite := &models.Invite{
		Code:      code,
		Email:     email,
		Used:      false,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // Expira em 7 dias
		CreatedAt: time.Now(),
	}

	err = s.userRepo.InsertInvite(ctx, invite)
	if err != nil {
		return nil, err
	}

	return invite, nil
}

func (s *UserService) Register(ctx context.Context, email, password, inviteCode string) error {
	// Verificar se o convite existe e é válido
	invite, err := s.userRepo.FindInviteByCode(ctx, inviteCode)
	if err != nil {
		return errors.New("convite inválido ou expirado")
	}
	if invite.Email != email || invite.Used || invite.ExpiresAt.Before(time.Now()) {
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

	err = s.userRepo.Insert(ctx, user)
	if err != nil {
		return err
	}

	// Marcar convite como usado
	invite.Used = true
	err = s.userRepo.UpdateInvite(ctx, invite)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, email string, password string) (string, *models.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
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
