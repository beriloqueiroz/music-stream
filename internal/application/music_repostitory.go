package application

import (
	"context"

	domain "github.com/beriloqueiroz/music-stream/internal/domain/entities"
)

type SearchResult struct {
	MusicList []*domain.Music
	Total     int
}

type MusicRepository interface {
	FindByID(ctx context.Context, id string) (*domain.Music, error)
	Create(ctx context.Context, music *domain.Music) (string, error)
	Search(ctx context.Context, query string, page int, limit int) (*SearchResult, error)
}
