package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	pb "github.com/beriloqueiroz/music-stream/api/proto"
	"github.com/beriloqueiroz/music-stream/internal/auth"
	"github.com/beriloqueiroz/music-stream/internal/music"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

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

	// Configuração do serviço de música
	musicService := music.NewMusicService(db, "./storage/music")

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
	mux.HandleFunc("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		authHandler.Register(w, r)
	})

	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		authHandler.Login(w, r)
	})

	// Rota protegida para criar convites (apenas admin)
	mux.Handle("/api/invites", authService.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		authHandler.CreateInvite(w, r)
	})))

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
