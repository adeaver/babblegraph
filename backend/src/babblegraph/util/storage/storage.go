package storage

// TODO: refactor this package
/*

type storage interface {
    Read(directory str, fileName str) []byte, err
    Write(directory str, file File)
}

type File struct {
    fileName str
    data []byte
    contentType ContentType
    accessControl AccessContorl
}

func CreateFileWithAccessControl(fileName string, accessControl) (*file, error)
func (f *File) WriteToFile(data []byte)

type remoteStorage struct {
    s3Storage util.s3.storage
}

func (r remoteStorage) Write
func (r remoteStorage) Read

var RemoteStorage

do local storage as well

// util.s3

type S3Client struct {
    Credentials
    BucketName (prod-spaces-1)
}

GetS3ClientForEnvironment

*/

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

type UploadDataInput struct {
	BucketName  string
	FileName    string
	Data        string
	ContentType ContentType
	IsPublic    bool
}

func (s *S3Storage) UploadData(input UploadDataInput) error {
	var acl *string
	if input.IsPublic {
		acl = ptr.String("public-read")
	}
	_, err := s.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(input.BucketName),
		Key:         aws.String(input.FileName),
		Body:        strings.NewReader(input.Data),
		ContentType: aws.String(input.ContentType.Str()),
		ACL:         acl,
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
