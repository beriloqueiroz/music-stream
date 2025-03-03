package application

import (
	"context"

	"github.com/beriloqueiroz/music-stream/pkg/models"
)

type PlaylistRepository interface {
	Create(ctx context.Context, playlist *models.Playlist) (string, error)
	FindByID(ctx context.Context, id string) (*models.Playlist, error)
	FindByUserID(ctx context.Context, userID string) ([]*models.Playlist, error)
	Update(ctx context.Context, playlist *models.Playlist) error
	Delete(ctx context.Context, id string) error
}
