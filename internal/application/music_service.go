package application

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"time"

	pb "github.com/beriloqueiroz/music-stream/api/proto"
	domain "github.com/beriloqueiroz/music-stream/internal/domain/entities"
	"github.com/beriloqueiroz/music-stream/pkg/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MusicService struct {
	pb.UnimplementedMusicServiceServer
	storage   storage.Storage
	musicRepo MusicRepository
}

func NewMusicService(storage storage.Storage, musicRepo MusicRepository) *MusicService {
	return &MusicService{
		storage:   storage,
		musicRepo: musicRepo,
	}
}

func (s *MusicService) GetMusic(ctx context.Context, id string) (*domain.Music, error) {
	music, err := s.musicRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return music, nil
}

func (s *MusicService) StreamMusic(req *pb.StreamRequest, stream pb.MusicService_StreamMusicServer) error {
	log.Println("StreamMusic chamado com query:", req)

	ctx := stream.Context()
	music, err := s.GetMusic(ctx, req.MusicId)
	if err != nil {
		return err
	}
	reader, err := s.storage.GetItem(music.StorageID.Hex())
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
			Type:           music.Type,
		}

		if err := stream.Send(chunk); err != nil {
			return err
		}
		sequence++
	}

	return nil
}

func (s *MusicService) UploadMusic(stream pb.MusicService_UploadMusicServer) error {
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

	if metadata.Title == "" {
		return errors.New("título não fornecido")
	}

	if metadata.Type == "" {
		return errors.New("tipo de arquivo não fornecido")
	}

	storageId := primitive.NewObjectID()

	// Salvar no storage
	err := s.storage.SaveItem(storageId.Hex(), &buffer)
	if err != nil {
		return err
	}

	musicID := primitive.NewObjectID()

	// salve AlbumArt image and get Id
	storageAlbumImageId := primitive.NewObjectID()

	// Salvar no storage
	if metadata.AlbumArt != nil {
		ext := ".jpeg"
		if metadata.AlbumArtType != "" {
			ext = "." + metadata.AlbumArtType
		}
		bufferAlbumArt := bytes.NewBuffer(metadata.AlbumArt)
		name := storageAlbumImageId.Hex() + ext
		err = s.storage.SaveItem(name, bufferAlbumArt)
		if err != nil {
			return err
		}
	}

	// Salvar no MongoDB
	music := &domain.Music{
		ID:        musicID,
		Title:     metadata.Title,
		Artist:    metadata.Artist,
		Album:     metadata.Album,
		StorageID: storageId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Type:      metadata.Type,
		Metadata: &domain.MusicMetadata{
			Title:    metadata.Title,
			Artist:   metadata.Artist,
			Album:    metadata.Album,
			Type:     metadata.Type,
			Year:     metadata.Year,
			Genre:    metadata.Genre,
			Composer: metadata.Composer,
			Label:    metadata.Label,
			AlbumArt: storageAlbumImageId.Hex(),
			Comments: metadata.Comments,
			Isrc:     metadata.Isrc,
			Url:      metadata.Url,
		},
	}

	id, err := s.musicRepo.Create(stream.Context(), music)
	if err != nil {
		// Se falhar, tenta remover do storage
		_ = s.storage.DeleteItem(storageId.Hex())
		return err
	}

	return stream.SendAndClose(&pb.UploadResponse{
		MusicId: id,
		Success: true,
		Message: "Música enviada com sucesso",
	})
}

func (s *MusicService) SearchMusic(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Println("SearchMusic chamado com query:", in.Query, "page:", in.Page, "limit:", in.PageSize)
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
			Id:     music.ID.Hex(),
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
