package blog

import (
	"babblegraph/util/env"
	"babblegraph/util/storage"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Content []ContentNode

type ContentNode struct {
	Type ContentNodeType `json:"type"`
	Body interface{}     `json:"body"`
}

type ContentNodeType string

const (
	ContentNodeTypeHeading   ContentNodeType = "heading"
	ContentNodeTypeParagraph ContentNodeType = "paragraph"
)

type Heading struct {
	Text string `json:"text"`
}

type Paragraph struct {
	Text string `json:"text"`
}

func getContentFileNameforID(id id) string {
	return fmt.Sprintf("blog-content/%s/content/%s.json", env.MustEnvironmentName(), id)
}

func GetContentForBlog(tx *sqlx.Tx, s3Storage *storage.S3Storage, urlPath string) ([]ContentNode, error) {
	metadata, err := getBlogPostMetadataByURLPath(tx, urlPath)
	if err != nil {
		return nil, err
	}
	contentStr, err := s3Storage.GetData("prod-spaces-1", getContentFileNameforID(metadata.ID))
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchKey") {
			return nil, nil
		}
		return nil, err
	}
	var content []ContentNode
	if err := json.Unmarshal([]byte(*contentStr), &content); err != nil {
		return nil, err
	}
	return content, nil
}

func UpsertContentForBlog(tx *sqlx.Tx, s3Storage *storage.S3Storage, urlPath string, content []ContentNode) error {
	if err := verifyContent(content); err != nil {
		return err
	}
	metadata, err := getBlogPostMetadataByURLPath(tx, urlPath)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(content)
	if err != nil {
		return err
	}
	return s3Storage.UploadData(storage.UploadDataInput{
		// TODO(staging environment): make this dynamic
		BucketName:  "prod-spaces-1",
		FileName:    getContentFileNameforID(metadata.ID),
		Data:        string(bytes),
		ContentType: storage.ContentTypeApplicationJSON,
	})
}

func verifyContent(content []ContentNode) error {
	var errs []string
	for idx, node := range content {
		switch node.Type {
		case ContentNodeTypeHeading:
			_, ok := node.Body.(Heading)
			if !ok {
				errs = append(errs, fmt.Sprintf("Node %d has type heading, but the body does not match", idx))
			}
		case ContentNodeTypeParagraph:
			_, ok := node.Body.(Paragraph)
			if !ok {
				errs = append(errs, fmt.Sprintf("Node %d has type paragraph, but the body does not match", idx))
			}
		default:
			errs = append(errs, fmt.Sprintf("Node %d has unrecognized type %s", idx, node.Type))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "\n"))
	}
	return nil
}
