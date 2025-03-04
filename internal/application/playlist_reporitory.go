package application

import (
	"context"

	domain "github.com/beriloqueiroz/music-stream/internal/domain/entities"
)

type PlaylistRepository interface {
	Create(ctx context.Context, playlist *domain.Playlist) (string, error)
	FindByID(ctx context.Context, id string, ownerID string) (*domain.Playlist, error)
	List(ctx context.Context, ownerID string, page int, limit int) ([]*domain.Playlist, error)
	Update(ctx context.Context, playlist *domain.Playlist) error
	Delete(ctx context.Context, id string, ownerID string) error
}
