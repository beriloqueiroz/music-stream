package grpc_server

import (
	"log"
	"net"

	pb "github.com/beriloqueiroz/music-stream/api/proto"
	"github.com/beriloqueiroz/music-stream/internal/music"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	grpcServer   *grpc.Server
	musicService *music.Service
	db           *mongo.Database
}

func NewGrpcServer(musicService *music.Service, db *mongo.Database) *GrpcServer {
	grpcServer := grpc.NewServer()
	pb.RegisterMusicServiceServer(grpcServer, musicService)

	return &GrpcServer{
		grpcServer:   grpcServer,
		musicService: musicService,
		db:           db,
	}
}

func (s *GrpcServer) Start() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Falha ao iniciar servidor gRPC: %v", err)
	}
	log.Printf("Servidor gRPC iniciado na porta 50051")
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("Falha ao servir gRPC: %v", err)
	}
}
