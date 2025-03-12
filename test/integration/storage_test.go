package integration

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	s3storage "github.com/beriloqueiroz/music-stream/pkg/storage/s3"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestStorageIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Pulando testes de integração")
	}

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "minio/minio",
		ExposedPorts: []string{"9000/tcp"},
		Env: map[string]string{
			"MINIO_ROOT_USER":     "minioadmin",
			"MINIO_ROOT_PASSWORD": "minioadmin",
			"MINIO_REGION":        "us-east-1",
		},
		Cmd: []string{"server", "/data"},
		WaitingFor: wait.ForAll(
			wait.ForLog("MinIO Object Storage Server"),
			wait.ForListeningPort("9000/tcp"),
		),
	}

	minioContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer minioContainer.Terminate(ctx)

	time.Sleep(2 * time.Second)

	mappedPort, err := minioContainer.MappedPort(ctx, "9000")
	if err != nil {
		t.Fatal(err)
	}

	endpoint := "http://localhost:" + mappedPort.Port()
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(endpoint),
		Region:           aws.String("us-east-1"),
		Credentials:      credentials.NewStaticCredentials("minioadmin", "minioadmin", ""),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})
	assert.NoError(t, err)

	s3Client := s3.New(sess)
	_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String("test-bucket"),
	})
	assert.NoError(t, err)

	storage := s3storage.NewS3Storage("test-bucket", sess)

	testData := []byte("test music data")
	err = storage.SaveItem("test-id", bytes.NewReader(testData))
	assert.NoError(t, err)

	reader, err := storage.GetItem("test-id")
	assert.NoError(t, err)
	defer reader.Close()

	data, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, testData, data)

	err = storage.DeleteItem("test-id")
	assert.NoError(t, err)
}
