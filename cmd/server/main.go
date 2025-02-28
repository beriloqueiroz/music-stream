package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	pb "github.com/beriloqueiroz/music-stream/api/proto"
	"github.com/beriloqueiroz/music-stream/internal/auth"
	"github.com/beriloqueiroz/music-stream/internal/music"
	"github.com/beriloqueiroz/music-stream/internal/playlist"
	"github.com/beriloqueiroz/music-stream/pkg/storage"
	"github.com/beriloqueiroz/music-stream/pkg/storage/s3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func getStorage() storage.MusicStorage {
	if os.Getenv("ENV") == "production" {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
		})
		if err != nil {
			log.Fatal(err)
		}
		return s3.NewS3Storage(os.Getenv("S3_BUCKET"), sess)
	}

	// Desenvolvimento: usar MinIO
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(os.Getenv("MINIO_ENDPOINT")),
		Region:   aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("MINIO_ACCESS_KEY"),
			os.Getenv("MINIO_SECRET_KEY"),
			"",
		),
		S3ForcePathStyle:              aws.Bool(true),
		DisableSSL:                    aws.Bool(true),
		S3DisableContentMD5Validation: aws.Bool(true),
		DisableEndpointHostPrefix:     aws.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}
	return s3.NewS3Storage(os.Getenv("BUCKET_NAME"), sess)
}

func main() {
	// Configuração do MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "musicstream"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Ping no MongoDB para verificar conexão
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(dbName)

	// Configuração do serviço de autenticação
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "bb#123$joao"
	}

	authService := auth.NewAuthService(db, jwtSecret)
	authHandler := auth.NewHandler(authService)

	// Configurar S3 (exemplo)
	storage := getStorage()
	musicService := music.NewMusicService(db, storage)

	// Configuração do gRPC
	grpcServer := grpc.NewServer()
	pb.RegisterMusicServiceServer(grpcServer, musicService)

	// Iniciar servidor gRPC
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Falha ao iniciar servidor gRPC: %v", err)
		}
		log.Printf("Servidor gRPC iniciado na porta 50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Falha ao servir gRPC: %v", err)
		}
	}()

	// Configuração das rotas
	mux := http.NewServeMux()

	// Rotas de autenticação
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/invites", authHandler.CreateInvite)

	// Rotas de playlists

	playlistService := playlist.NewPlaylistService(db)
	playlistHandler := playlist.NewHandler(playlistService)

	mux.Handle("POST /api/playlists", authService.AuthMiddleware(http.HandlerFunc(playlistHandler.CreatePlaylist)))
	mux.Handle("GET /api/playlists/{id}", authService.AuthMiddleware(http.HandlerFunc(playlistHandler.GetPlaylist)))
	mux.Handle("PUT /api/playlists/{id}", authService.AuthMiddleware(http.HandlerFunc(playlistHandler.UpdatePlaylist)))
	mux.Handle("DELETE /api/playlists/{id}", authService.AuthMiddleware(http.HandlerFunc(playlistHandler.DeletePlaylist)))
	mux.Handle("POST /api/playlists/{id}/musics", authService.AuthMiddleware(http.HandlerFunc(playlistHandler.AddMusicInPlaylist)))
	mux.Handle("DELETE /api/playlists/{id}/musics/{musicId}", authService.AuthMiddleware(http.HandlerFunc(playlistHandler.RemoveMusicInPlaylist)))
	mux.Handle("GET /api/playlists/{id}/musics", authService.AuthMiddleware(http.HandlerFunc(playlistHandler.GetPlaylist)))

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
