package music

import (
	"context"
	"io"

	pb "github.com/beriloqueiroz/music-stream/api/proto"
	"github.com/beriloqueiroz/music-stream/pkg/models"
	"github.com/beriloqueiroz/music-stream/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
