package rest_server

import (
	"log"
	"net/http"
	"time"

	"github.com/beriloqueiroz/music-stream/internal/application"
	rest_server_playlist "github.com/beriloqueiroz/music-stream/internal/infra/rest/playlist"
	rest_server_user "github.com/beriloqueiroz/music-stream/internal/infra/rest/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type RestServer struct {
	db *mongo.Database
}

func NewRestServer(db *mongo.Database) *RestServer {
	return &RestServer{db: db}
}

func (s *RestServer) Start(jwtSecret string, userRepo application.UserRepository) {
	// Configuração das rotas
	mux := http.NewServeMux()

	authService := application.NewUserService(userRepo, []byte(jwtSecret))
	authHandler := rest_server_user.NewUserHandler(authService)
	authMiddlewares := rest_server_user.NewUserMiddlewares([]byte(jwtSecret))
	// Rotas de autenticação
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.Handle("POST /api/invites", authMiddlewares.AuthMiddleware(http.HandlerFunc(authHandler.CreateInvite)))
	// Rotas de playlists
	playlistService := application.NewPlaylistService(s.db)
	playlistHandler := rest_server_playlist.NewPlaylistHandler(playlistService)

	mux.Handle("DELETE /api/playlists/{id}/musics/{musicId}", authMiddlewares.AuthMiddleware(http.HandlerFunc(playlistHandler.RemoveMusicInPlaylist)))
	mux.Handle("POST /api/playlists/{id}/musics", authMiddlewares.AuthMiddleware(http.HandlerFunc(playlistHandler.AddMusicInPlaylist)))
	mux.Handle("GET /api/playlists/{id}/musics", authMiddlewares.AuthMiddleware(http.HandlerFunc(playlistHandler.GetPlaylist)))
	mux.Handle("DELETE /api/playlists/{id}", authMiddlewares.AuthMiddleware(http.HandlerFunc(playlistHandler.DeletePlaylist)))
	mux.Handle("PUT /api/playlists/{id}", authMiddlewares.AuthMiddleware(http.HandlerFunc(playlistHandler.UpdatePlaylist)))
	mux.Handle("GET /api/playlists", authMiddlewares.AuthMiddleware(http.HandlerFunc(playlistHandler.GetPlaylists)))
	mux.Handle("POST /api/playlists", authMiddlewares.AuthMiddleware(http.HandlerFunc(playlistHandler.CreatePlaylist)))

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
