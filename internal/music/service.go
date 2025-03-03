package music

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	pb "github.com/beriloqueiroz/music-stream/api/proto"
	"github.com/beriloqueiroz/music-stream/internal/application"
	"github.com/beriloqueiroz/music-stream/pkg/models"
	"github.com/beriloqueiroz/music-stream/pkg/storage"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	pb.UnimplementedMusicServiceServer
	db        *mongo.Database
	storage   storage.MusicStorage
	musicRepo application.MusicRepository
}

func NewMusicService(db *mongo.Database, storage storage.MusicStorage, musicRepo application.MusicRepository) *Service {
	return &Service{
		db:        db,
		storage:   storage,
		musicRepo: musicRepo,
	}
}

func (s *Service) GetMusic(ctx context.Context, id string) (*models.Music, error) {
	music, err := s.musicRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return music, nil
}

func (s *Service) StreamMusic(req *pb.StreamRequest, stream pb.MusicService_StreamMusicServer) error {
	ctx := stream.Context()
	music, err := s.GetMusic(ctx, req.MusicId)
	if err != nil {
		return err
	}
	storageID := music.StorageID
	reader, err := s.storage.GetMusic(storageID)
	if err != nil {
		return err
	}
	defer reader.Close()

	buffer := make([]byte, 1024*32)
	sequence := int32(0)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		chunk := &pb.AudioChunk{
			Data:           buffer[:n],
			SequenceNumber: sequence,
		}

		if err := stream.Send(chunk); err != nil {
			return err
		}
		sequence++
	}

	return nil
}

func (s *Service) UploadMusic(stream pb.MusicService_UploadMusicServer) error {
	var metadata *pb.MusicMetadata

	// Buffer para acumular os chunks
	var buffer bytes.Buffer

	// Receber chunks do cliente
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Primeiro chunk contém metadata
		switch data := req.Data.(type) {
		case *pb.UploadRequest_Metadata:
			metadata = data.Metadata
		case *pb.UploadRequest_ChunkData:
			_, err = buffer.Write(data.ChunkData)
			if err != nil {
				return err
			}
		}
	}

	if metadata == nil {
		return errors.New("metadata não fornecida")
	}

	storageUUID := uuid.New().String()

	// Salvar no storage
	err := s.storage.SaveMusic(storageUUID, &buffer)
	if err != nil {
		return err
	}

	// Salvar no MongoDB
	music := &models.Music{
		Title:     metadata.Title,
		Artist:    metadata.Artist,
		Album:     metadata.Album,
		StorageID: storageUUID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := s.musicRepo.Create(stream.Context(), music)
	if err != nil {
		// Se falhar, tenta remover do storage
		_ = s.storage.DeleteMusic(storageUUID)
		return err
	}

	return stream.SendAndClose(&pb.UploadResponse{
		MusicId: id,
		Success: true,
		Message: "Música enviada com sucesso",
	})
}

func (s *Service) SearchMusic(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	query := in.Query
	page := in.Page
	limit := in.PageSize

	musics, err := s.musicRepo.Search(ctx, query, int(page), int(limit))
	if err != nil {
		return nil, err
	}

	var musicsList []*pb.Music
	for _, music := range musics.MusicList {
		musicsList = append(musicsList, &pb.Music{
			Id:     music.ID,
			Title:  music.Title,
			Artist: music.Artist,
			Album:  music.Album,
		})
	}

	return &pb.SearchResponse{
		MusicList: musicsList,
		Total:     int32(len(musicsList)),
	}, nil
}
