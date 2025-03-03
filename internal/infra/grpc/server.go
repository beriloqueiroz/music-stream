package grpc_server

import (
	"log"
	"net"

	pb "github.com/beriloqueiroz/music-stream/api/proto"
	"github.com/beriloqueiroz/music-stream/internal/application"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	grpcServer   *grpc.Server
	musicService *application.MusicService
}

func NewGrpcServer(musicService *application.MusicService) *GrpcServer {
	grpcServer := grpc.NewServer()
	pb.RegisterMusicServiceServer(grpcServer, musicService)

	return &GrpcServer{
		grpcServer:   grpcServer,
		musicService: musicService,
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
