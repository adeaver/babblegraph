package storage

import (
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"bytes"
	"io"
	"strings"

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

type ContentType string

const (
	ContentTypeApplicationJSON = "application/json"
)

func (c ContentType) Str() string {
	return string(c)
}

type UploadDataInput struct {
	BucketName  string
	FileName    string
	Data        string
	ContentType ContentType
}

func (s *S3Storage) UploadData(input UploadDataInput) error {
	_, err := s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(input.BucketName),
		Key:         aws.String(input.FileName),
		Body:        strings.NewReader(input.Data),
		ContentType: aws.String(input.ContentType.Str()),
	})
	return err
}

func (s *S3Storage) DeleteData(bucketName, fileName string) error {
	if _, err := s.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	}); err != nil {
		return err
	}
	return nil
}