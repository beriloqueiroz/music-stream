package s3

import (
	"bytes"
	"io"
	"testing"

	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockS3Client struct {
	s3iface.S3API
	mock.Mock
}

func (m *mockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	args := m.Called(input)
	if output := args.Get(0); output != nil {
		return output.(*s3.GetObjectOutput), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	args := m.Called(input)
	if output := args.Get(0); output != nil {
		return output.(*s3.PutObjectOutput), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockS3Client) DeleteObject(input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	args := m.Called(input)
	if output := args.Get(0); output != nil {
		return output.(*s3.DeleteObjectOutput), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockS3Client) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	args := m.Called(input)
	if output := args.Get(0); output != nil {
		return output.(*s3manager.UploadOutput), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestS3Storage_GetMusic(t *testing.T) {
	mockClient := new(mockS3Client)
	storage := &S3Storage{
		bucket: "test-bucket",
		client: mockClient,
	}

	testData := []byte("test music data")
	mockClient.On("GetObject", mock.Anything).Return(&s3.GetObjectOutput{
		Body: io.NopCloser(bytes.NewReader(testData)),
	}, nil)

	reader, err := storage.GetMusic("test-id")
	assert.NoError(t, err)

	data, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, testData, data)
}

func TestS3Storage_SaveMusic_SmallFile(t *testing.T) {
	mockClient := new(mockS3Client)
	storage := &S3Storage{
		bucket:        "test-bucket",
		client:        mockClient,
		maxMemorySize: 32 * 1024 * 1024,
	}

	testData := []byte("test music data")
	mockClient.On("PutObject", mock.MatchedBy(func(input *s3.PutObjectInput) bool {
		return *input.Bucket == "test-bucket" &&
			*input.Key == "test-id" &&
			*input.ContentLength == int64(len(testData))
	})).Return(&s3.PutObjectOutput{}, nil)

	err := storage.SaveMusic("test-id", bytes.NewReader(testData))
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestS3Storage_SaveMusic_LargeFile(t *testing.T) {
	mockClient := new(mockS3Client)
	storage := &S3Storage{
		bucket:        "test-bucket",
		client:        mockClient,
		maxMemorySize: 10, // Forçar uso do uploader
	}

	// Criar dados maiores que maxMemorySize
	testData := make([]byte, 20)
	for i := range testData {
		testData[i] = byte(i)
	}

	// Mock para o uploader
	mockClient.On("Upload", mock.MatchedBy(func(input *s3manager.UploadInput) bool {
		return *input.Bucket == "test-bucket" && *input.Key == "test-id"
	})).Return(&s3manager.UploadOutput{}, nil)

	err := storage.SaveMusic("test-id", bytes.NewReader(testData))
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestS3Storage_DeleteMusic(t *testing.T) {
	mockClient := new(mockS3Client)
	storage := &S3Storage{
		bucket: "test-bucket",
		client: mockClient,
	}

	mockClient.On("DeleteObject", &s3.DeleteObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("test-id"),
	}).Return(&s3.DeleteObjectOutput{}, nil)

	err := storage.DeleteMusic("test-id")
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestS3Storage_Errors(t *testing.T) {
	mockClient := new(mockS3Client)
	storage := &S3Storage{
		bucket: "test-bucket",
		client: mockClient,
	}

	expectedErr := errors.New("s3 error")

	// Teste de erro no GetMusic
	mockClient.On("GetObject", mock.Anything).Return(nil, expectedErr)
	_, err := storage.GetMusic("test-id")
	assert.Error(t, err)

	// Teste de erro no SaveMusic
	mockClient.On("PutObject", mock.Anything).Return(nil, expectedErr)
	err = storage.SaveMusic("test-id", bytes.NewReader([]byte("data")))
	assert.Error(t, err)

	// Teste de erro no DeleteMusic
	mockClient.On("DeleteObject", mock.Anything).Return(nil, expectedErr)
	err = storage.DeleteMusic("test-id")
	assert.Error(t, err)
}

// ... mais testes para SaveMusic e DeleteMusic
