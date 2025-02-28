package s3

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Storage struct {
	bucket string
	client *s3.S3
}

func NewS3Storage(bucket string, client *s3.S3) *S3Storage {
	return &S3Storage{
		bucket: bucket,
		client: client,
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
	input := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(id),
		Body:   aws.ReadSeekCloser(data),
	}

	_, err := s.client.PutObject(input)
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
