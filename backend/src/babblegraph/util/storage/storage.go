package storage

import (
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Storage struct {
	s3Client *s3.S3
}

func NewS3StorageForEnvironment() *S3Storage {
	key := env.MustEnvironmentVariable("S3_KEY")
	secret := env.MustEnvironmentVariable("S3_SECRET")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String("https://nyc3.digitaloceanspaces.com"),
		Region:      aws.String("us-east-1"),
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)
	return &S3Storage{
		s3Client: s3Client,
	}
}

func (s *S3Storage) GetData(bucketName, fileName string) (*string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	}

	result, err := s.s3Client.GetObject(input)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	_, err = io.Copy(&b, result.Body)
	if err != nil {
		return nil, err
	}
	return ptr.String(b.String()), nil
}
