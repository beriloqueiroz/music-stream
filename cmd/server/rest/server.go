package rest_server

import (
	"log"
	"net/http"
	"time"

	"github.com/beriloqueiroz/music-stream/internal/auth"
	"github.com/beriloqueiroz/music-stream/internal/playlist"
	"go.mongodb.org/mongo-driver/mongo"
)

type RestServer struct {
	db *mongo.Database
}

func NewRestServer(db *mongo.Database) *RestServer {
	return &RestServer{db: db}
}

func (s *RestServer) Start(jwtSecret string) {
	// Configuração das rotas
	mux := http.NewServeMux()

	authService := auth.NewAuthService(s.db, jwtSecret)
	authHandler := auth.NewHandler(authService)
	// Rotas de autenticação
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.Handle("POST /api/invites", authService.AuthMiddleware(http.HandlerFunc(authHandler.CreateInvite)))
	// Rotas de playlists
	playlistService := playlist.NewPlaylistService(s.db)
	playlistHandler := playlist.NewHandler(playlistService)

	mux.Handle("DELETE /api/playlists/{id}/musics/{musicId}", authService.AuthMiddleware(http.HandlerFunc(playlistHandler.RemoveMusicInPlaylist)))
	mux.Handle("POST /api/playlists/{id}/musics", authService.AuthMiddleware(http.HandlerFunc(playlistHandler.AddMusicInPlaylist)))
	mux.Handle("GET /api/playlists/{id}/musics", authService.AuthMiddleware(http.HandlerFunc(playlistHandler.GetPlaylist)))
	mux.Handle("DELETE /api/playlists/{id}", authService.AuthMiddleware(http.HandlerFunc(playlistHandler.DeletePlaylist)))
	mux.Handle("PUT /api/playlists/{id}", authService.AuthMiddleware(http.HandlerFunc(playlistHandler.UpdatePlaylist)))
	mux.Handle("GET /api/playlists", authService.AuthMiddleware(http.HandlerFunc(playlistHandler.GetPlaylists)))
	mux.Handle("POST /api/playlists", authService.AuthMiddleware(http.HandlerFunc(playlistHandler.CreatePlaylist)))

	// Configuração do servidor
	srv := &http.Server{
		Handler:      mux,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Servidor iniciado na porta 8080")
	log.Fatal(srv.ListenAndServe())
}
