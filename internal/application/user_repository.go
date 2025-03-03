package application

import (
	"context"

	"github.com/beriloqueiroz/music-stream/pkg/models"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*models.User, error)
	Insert(ctx context.Context, user *models.User) error
	InsertInvite(ctx context.Context, invite *models.Invite) error
	UpdateInvite(ctx context.Context, invite *models.Invite) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindInviteByCode(ctx context.Context, code string) (*models.Invite, error)
}
