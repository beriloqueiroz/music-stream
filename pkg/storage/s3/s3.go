package s3

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Client interface {
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	DeleteObject(*s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error)
	AbortMultipartUpload(*s3.AbortMultipartUploadInput) (*s3.AbortMultipartUploadOutput, error)
	CreateMultipartUpload(*s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error)
	CompleteMultipartUpload(*s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error)
	UploadPart(*s3.UploadPartInput) (*s3.UploadPartOutput, error)
}

type S3Storage struct {
	bucket string
	client s3iface.S3API
	// Tamanho máximo para upload em memória (ex: 32MB)
	maxMemorySize int64
}

func NewS3Storage(bucket string, sess *session.Session) *S3Storage {
	return &S3Storage{
		bucket:        bucket,
		client:        s3.New(sess),
		maxMemorySize: 32 * 1024 * 1024, // 32MB
	}
}

func (s *S3Storage) GetMusic(id string) (io.ReadCloser, error) {
	output, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(id),
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}

func (s *S3Storage) SaveMusic(id string, data io.Reader) error {
	// Primeiro, tenta ler um pouco do arquivo para ver o tamanho
	buf := &bytes.Buffer{}
	n, err := io.CopyN(buf, data, s.maxMemorySize+1)
	if err != nil && err != io.EOF {
		return err
	}

	// Se o arquivo é menor que maxMemorySize, usa o método simples
	if n <= s.maxMemorySize {
		input := &s3.PutObjectInput{
			Bucket:        aws.String(s.bucket),
			Key:           aws.String(id),
			Body:          bytes.NewReader(buf.Bytes()),
			ContentLength: aws.Int64(n),
		}
		_, err := s.client.PutObject(input)
		return err
	}

	// Para arquivos grandes, usa o uploader
	uploader := s3manager.NewUploaderWithClient(s.client)

	// Combina o buffer já lido com o resto do reader
	multiReader := io.MultiReader(buf, data)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(id),
		Body:   multiReader,
	})
	return err
}

func (s *S3Storage) DeleteMusic(id string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(id),
	}

	_, err := s.client.DeleteObject(input)
	return err
}

// ... implementar outros métodos ...
