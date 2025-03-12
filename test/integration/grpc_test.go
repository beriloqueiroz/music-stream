package integration

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	pb "github.com/beriloqueiroz/music-stream/api/proto"
	"github.com/beriloqueiroz/music-stream/internal/application"
	grpc_server "github.com/beriloqueiroz/music-stream/internal/infra/grpc"
	"github.com/beriloqueiroz/music-stream/internal/infra/mongodb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPCIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando testes de integração")
	}

	ctx := context.Background()

	testHelper := NewTestHelper()

	container, err := testHelper.StartMongoDBContainer(ctx)
	if err != nil {
		log.Fatalf("Erro ao iniciar container: %v", err)
		os.Exit(1)
	}
	defer container.Terminate(ctx)

	database, err := testHelper.ConnectToMongoDB(ctx, container)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
		os.Exit(1)
	}
	defer database.Disconnect(ctx)

	musicRepo := mongodb.NewMongoMusicRepository(database.Database("music-stream"))

	minioContainer, err := testHelper.StartMinioContainer(ctx)
	if err != nil {
		log.Fatalf("Erro ao iniciar container Minio: %v", err)
		os.Exit(1)
	}
	defer minioContainer.Terminate(ctx)

	storage, err := testHelper.GetMinioStorage(ctx, minioContainer)
	if err != nil {
		log.Fatalf("Erro ao obter storage Minio: %v", err)
		os.Exit(1)
	}
	musicService := application.NewMusicService(storage, musicRepo)

	grpcServer := grpc_server.NewGrpcServer(musicService, "50052")
	go grpcServer.Start()

	time.Sleep(2 * time.Second)

	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Erro ao conectar ao gRPC: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewMusicServiceClient(conn)

	t.Run("TestUploadMusic", func(t *testing.T) {
		// test create music
		stream, err := client.UploadMusic(ctx)
		if err != nil {
			log.Fatalf("Erro ao criar stream: %v", err)
			os.Exit(1)
		}
		err = stream.Send(&pb.UploadRequest{
			Data: &pb.UploadRequest_Metadata{
				Metadata: &pb.MusicMetadata{
					Title:    "Yesterday",
					Artist:   "The Beatles",
					Album:    "The Beatles 1967-1970",
					Type:     "mp3",
					Year:     1967,
					Genre:    "Rock",
					Composer: "John Lennon",
					Label:    "Apple Records",
					AlbumArt: []byte{},
					Comments: "This is a test comment",
					Isrc:     "ABCD12345678",
					Url:      "https://example.com/yesterday",
				},
			},
		})
		if err != nil {
			log.Fatalf("Erro ao enviar metadata: %v", err)
			os.Exit(1)
		}

		//l list files in test_data folder
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Erro ao obter diretório atual: %v", err)
			os.Exit(1)
		}
		file, err := os.Open(filepath.Join(currentDir, "test_data/yesterday.mp3"))
		if err != nil {
			log.Fatalf("Erro ao abrir arquivo: %v", err)
			os.Exit(1)
		}
		defer file.Close()

		// Enviar chunks do arquivo
		buffer := make([]byte, 1024*32)
		for {
			n, err := file.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Erro ao ler arquivo: %v", err)
				os.Exit(1)
			}

			err = stream.Send(&pb.UploadRequest{
				Data: &pb.UploadRequest_ChunkData{
					ChunkData: buffer[:n],
				},
			})
			if err != nil {
				log.Fatalf("Erro ao enviar chunk: %v", err)
				os.Exit(1)
			}
		}

		response, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Erro ao receber resposta: %v", err)
			os.Exit(1)
		}

		assert.NotNil(t, response)
	})

	musicID := ""

	t.Run("TestSearchMusic", func(t *testing.T) {
		// test search music
		response, err := client.SearchMusic(ctx, &pb.SearchRequest{
			Query: "Yesterday",
		})
		if err != nil {
			log.Fatalf("Erro ao buscar música: %v", err)
			os.Exit(1)
		}

		assert.NotNil(t, response)
		assert.Equal(t, response.MusicList[0].Title, "Yesterday")
		musicID = response.MusicList[0].Id
	})

	t.Run("TestSearchMusicWithInvalidQuery", func(t *testing.T) {
		// test search music with invalid query
		response, err := client.SearchMusic(ctx, &pb.SearchRequest{
			Query: "Invalid Query",
		})
		if err != nil {
			log.Fatalf("Erro ao buscar música: %v", err)
			os.Exit(1)
		}

		assert.NotNil(t, response)
		assert.Equal(t, len(response.MusicList), 0)
	})

	t.Run("TestPlayStreamMusic", func(t *testing.T) {
		// test get music
		responseStream, err := client.StreamMusic(ctx, &pb.StreamRequest{
			MusicId: musicID,
		})
		if err != nil {
			log.Fatalf("Erro ao buscar música: %v", err)
			os.Exit(1)
		}

		if err != nil {
			log.Fatal(err)
		}

		tempFile, err := os.CreateTemp("", "stream*.mp3")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(tempFile.Name())

		for {
			chunk, err := responseStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if _, err := tempFile.Write(chunk.Data); err != nil {
				log.Fatal(err)
			}
		}

		assert.NotNil(t, responseStream)
		assert.NotNil(t, tempFile)
		assert.FileExists(t, tempFile.Name())
	})
}
