package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Storage struct {
	bucket string
	client *s3.Client
}

func (s *S3Storage) GetMusic(id string) (io.ReadCloser, error) {
	output, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(id),
	})
	if err != nil {
		return nil, err
	}
	return output.Body, nil
}

// ... implementar outros métodos ...
