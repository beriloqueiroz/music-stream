package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/beriloqueiroz/music-stream/pkg/storage"
	"github.com/beriloqueiroz/music-stream/pkg/storage/s3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TestHelper struct {
}

func NewTestHelper() *TestHelper {
	return &TestHelper{}
}

func (h *TestHelper) StartMongoDBContainer(ctx context.Context) (testcontainers.Container, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	log.Println("Project directory:", filepath.Join(currentDir, "test_data/mongo-init.js"))

	containerRequest := testcontainers.ContainerRequest{
		Image: "mongo:latest",
		// Name:         "music-stream-test",
		ExposedPorts: []string{"27018:27018"},
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "root",
			"MONGO_INITDB_ROOT_PASSWORD": "root",
			"MONGO_INITDB_DATABASE":      "music-stream",
		},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      filepath.Join(currentDir, "test_data/mongo-init.js"),
				ContainerFilePath: "/docker-entrypoint-initdb.d/mongo-init.js",
				FileMode:          0o644,
			},
		},
		Cmd: []string{"mongod", "--auth", "--port", "27018"},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	log.Println("Container MongoDB iniciado com sucesso")

	return container, nil
}

func (h *TestHelper) ConnectToMongoDB(ctx context.Context, container testcontainers.Container) (*mongo.Client, error) {
	mappedPort, err := container.MappedPort(ctx, "27018")

	if err != nil {
		log.Fatalf("Erro ao obter porta mapeada: %v", err)
		os.Exit(1)
	}

	dbURL := fmt.Sprintf("mongodb://root:root@localhost:%s/music-stream?authSource=admin", mappedPort.Port())

	db, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURL))
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
		os.Exit(1)
	}

	err = db.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("MongoDB conectado com sucesso")

	return db, nil
}

// Erro ao iniciar container Minio: create container: Error response from daemon: pull access denied for minio, repository does not exist or may require 'docker login': denied: requested access to the resource is denied
func (h *TestHelper) StartMinioContainer(ctx context.Context) (testcontainers.Container, error) {
	net, err := network.New(ctx)
	if err != nil {
		return nil, err
	}
	containerRequest := testcontainers.ContainerRequest{
		Image:        "minio/minio",
		Networks:     []string{net.Name},
		Name:         "minio-test",
		ExposedPorts: []string{"8000:8000", "8001:8001"},
		Env: map[string]string{
			"MINIO_ACCESS_KEY":    "minioadmin",
			"MINIO_SECRET_KEY":    "minioadmin",
			"MINIO_ROOT_USER":     "minioadmin",
			"MINIO_ROOT_PASSWORD": "minioadmin",
		},
		Cmd: []string{"server", "/data", "--console-address", ":8001", "--address", ":8000"},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequest,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	log.Println("Container Minio iniciado com sucesso")

	containerRequestInit := testcontainers.ContainerRequest{
		Name:       "minio-init-test",
		Image:      "minio/mc",
		Networks:   []string{net.Name},
		Entrypoint: []string{"/bin/sh", "-c", "sleep 1 && mc alias set myminio http://minio-test:8000 minioadmin minioadmin && mc mb myminio/music-bucket-test"},
	}
	containerInit, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerRequestInit,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	log.Println("Container Minio Init iniciado com sucesso", containerInit.GetContainerID())

	return container, nil
}

func (h *TestHelper) GetMinioStorage(ctx context.Context, container testcontainers.Container) (storage.MusicStorage, error) {

	endpoint, err := container.Endpoint(ctx, "http")
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(endpoint),
		Region:   aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			"minioadmin",
			"minioadmin",
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
	return s3.NewS3Storage("music-bucket-test", sess), nil
}
