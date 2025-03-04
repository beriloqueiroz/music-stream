package application

import (
	"context"

	domain "github.com/beriloqueiroz/music-stream/internal/domain/entities"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Insert(ctx context.Context, user *domain.User) error
	InsertInvite(ctx context.Context, invite *domain.Invite) error
	UpdateInvite(ctx context.Context, invite *domain.Invite) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindInviteByCode(ctx context.Context, code string) (*domain.Invite, error)
}
