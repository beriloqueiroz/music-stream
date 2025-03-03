package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	grpc_server "github.com/beriloqueiroz/music-stream/cmd/server/grpc"
	rest_server "github.com/beriloqueiroz/music-stream/cmd/server/rest"
	"github.com/beriloqueiroz/music-stream/internal/infra/mongodb"
	"github.com/beriloqueiroz/music-stream/internal/music"
	"github.com/beriloqueiroz/music-stream/pkg/storage"
	"github.com/beriloqueiroz/music-stream/pkg/storage/s3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	client, dbName, err := connectMongoDB(ctx)
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

	musicRepo := mongodb.NewMongoMusicRepository(db)

	// Configurar S3 (exemplo)
	storage := getStorage()
	musicService := music.NewMusicService(db, storage, musicRepo)

	grpcServer := grpc_server.NewGrpcServer(musicService, db)
	go grpcServer.Start()

	RestServer := rest_server.NewRestServer(db)
	RestServer.Start(jwtSecret)
}

func connectMongoDB(ctx context.Context) (*mongo.Client, string, error) {
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

	return client, dbName, nil
}
