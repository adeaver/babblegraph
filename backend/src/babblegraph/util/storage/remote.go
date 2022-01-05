package storage

import (
	"babblegraph/util/env"
	"bytes"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var RemoteStorage *remoteStorage = nil

type remoteStorage struct {
	bucketName *string
	s3Client   *s3.S3
}

func getRemoteStorageForEnvironment() *remoteStorage {
	key := env.MustEnvironmentVariable("S3_KEY")
	secret := env.MustEnvironmentVariable("S3_SECRET")
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String("https://nyc3.digitaloceanspaces.com"),
		Region:      aws.String("us-east-1"),
	}
	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)
	return &remoteStorage{
		bucketName: aws.String("prod-spaces-1"),
		s3Client:   s3Client,
	}
}

func (r *remoteStorage) Read(directory string, fileName string) ([]byte, error) {
	result, err := r.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: r.bucketName,
		Key:    aws.String(fmt.Sprintf("%s/%s", directory, fileName)),
	})
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, result.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *remoteStorage) Write(directory string, file File) error {
	_, err := r.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      r.bucketName,
		Key:         aws.String(fmt.Sprintf("%s/%s", directory, file.name)),
		Body:        bytes.NewReader(file.data),
		ContentType: aws.String(file.contentType.Str()),
		ACL:         aws.String(file.accessControl.Str()),
	})
	return err
}

func (r *remoteStorage) Delete(directory string, fileName string) error {
	_, err := r.s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: r.bucketName,
		Key:    aws.String(fmt.Sprintf("%s/%s", directory, fileName)),
	})
	return err
}
