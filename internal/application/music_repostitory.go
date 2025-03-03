package application

import (
	"context"

	"github.com/beriloqueiroz/music-stream/pkg/models"
)

type SearchResult struct {
	MusicList []*models.Music
	Total     int
}

type MusicRepository interface {
	FindByID(ctx context.Context, id string) (*models.Music, error)
	Create(ctx context.Context, music *models.Music) (string, error)
	Search(ctx context.Context, query string, page int, limit int) (*SearchResult, error)
}
