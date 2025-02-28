package music

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	pb "github.com/beriloqueiroz/music-stream/api/proto"
	"github.com/beriloqueiroz/music-stream/pkg/models"
	"github.com/beriloqueiroz/music-stream/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	pb.UnimplementedMusicServiceServer
	db         *mongo.Database
	musicsColl *mongo.Collection
	storage    storage.MusicStorage
}

func NewMusicService(db *mongo.Database, storage storage.MusicStorage) *Service {
	return &Service{
		db:         db,
		musicsColl: db.Collection("musics"),
		storage:    storage,
	}
}

func (s *Service) GetMusic(ctx context.Context, id string) (*models.Music, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	music := &models.Music{}
	err = s.musicsColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(music)
	if err != nil {
		return nil, err
	}

	return music, nil
}

func (s *Service) StreamMusic(req *pb.StreamRequest, stream pb.MusicService_StreamMusicServer) error {
	reader, err := s.storage.GetMusic(req.MusicId)
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

	// Gerar um único ID para usar em ambos os lugares
	id := primitive.NewObjectID()
	storageID := id.Hex()

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

	// Salvar no storage
	err := s.storage.SaveMusic(storageID, &buffer)
	if err != nil {
		return err
	}

	// Salvar no MongoDB
	music := &models.Music{
		ID:        id,
		Title:     metadata.Title,
		Artist:    metadata.Artist,
		Album:     metadata.Album,
		StorageID: storageID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.musicsColl.InsertOne(stream.Context(), music)
	if err != nil {
		// Se falhar, tenta remover do storage
		_ = s.storage.DeleteMusic(storageID)
		return err
	}

	return stream.SendAndClose(&pb.UploadResponse{
		MusicId: storageID,
		Success: true,
		Message: "Música enviada com sucesso",
	})
}

func (s *Service) SearchMusic(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	query := in.Query
	page := in.Page
	limit := in.PageSize

	musics, err := s.musicsColl.Find(ctx, bson.M{
		"$text": bson.M{"$search": query},
	}, options.Find().SetSkip(int64(page*limit)).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}

	var musicsList []*pb.Music
	for musics.Next(ctx) {
		var music = &pb.Music{}
		if err := musics.Decode(&music); err != nil {
			return nil, err
		}
		musicsList = append(musicsList, music)
	}

	return &pb.SearchResponse{
		MusicList: musicsList,
		Total:     int32(len(musicsList)),
	}, nil
}
