package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/beriloqueiroz/music-stream/api/proto"
	"github.com/beriloqueiroz/music-stream/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Conectar ao MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Conectar ao gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Fazer upload da música
	grpcClient := pb.NewMusicServiceClient(conn)
	musicFile := "test_data/yesterday.mp3" // Coloque o arquivo MP3 aqui

	file, err := os.Open(musicFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stream, err := grpcClient.UploadMusic(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Enviar metadata
	err = stream.Send(&pb.UploadRequest{
		Data: &pb.UploadRequest_Metadata{
			Metadata: &pb.MusicMetadata{
				Title:  "Yesterday",
				Artist: "The Beatles",
				Album:  "Help!",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Enviar chunks
	buffer := make([]byte, 1024*32)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		err = stream.Send(&pb.UploadRequest{
			Data: &pb.UploadRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	// Atualizar no MongoDB com letra e cifra
	db := client.Database("musicstream")
	coll := db.Collection("musics")

	musicID, err := primitive.ObjectIDFromHex(response.MusicId)
	if err != nil {
		log.Fatal(err)
	}

	music := models.Music{
		ID:        musicID,
		Title:     "Yesterday",
		Artist:    "The Beatles",
		Album:     "Help!",
		StorageID: response.MusicId,
		Lyrics: &models.Lyrics{
			Text: "Yesterday, all my troubles seemed so far away...",
			Timing: []models.Segment{
				{Start: 0.0, End: 2.5, Content: "Yesterday,"},
				{Start: 2.5, End: 5.0, Content: "all my troubles seemed"},
				{Start: 5.0, End: 7.5, Content: "so far away"},
			},
			Language: "en",
		},
		Tablature: &models.Tablature{
			Content: "Am Dm G7 C",
			Timing: []models.Segment{
				{Start: 0.0, End: 2.5, Content: "Am"},
				{Start: 2.5, End: 5.0, Content: "Dm"},
				{Start: 5.0, End: 6.0, Content: "G7"},
				{Start: 6.0, End: 7.5, Content: "C"},
			},
			Format: "chordpro",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": musicID}
	update := bson.M{"$set": music}

	_, err = coll.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Música inserida/atualizada com ID: %s", music.ID.Hex())
}
