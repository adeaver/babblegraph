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
	ContentNodeTypeImage     ContentNodeType = "image"
	ContentNodeTypeLink      ContentNodeType = "link"
	ContentNodeTypeList      ContentNodeType = "list"
)

type Heading struct {
	Text string `json:"text"`
}

type Paragraph struct {
	Text string `json:"text"`
}

type Image struct {
	Path    string  `json:"path"`
	AltText string  `json:"alt_text"`
	Caption *string `json:"caption,omitempty"`
}

type Link struct {
	DestinationURL string `json:"destination_url"`
	Text           string `json:"text"`
}

type List struct {
	Items []string `json:"items"`
	Type  ListType `json:"type"`
}

type ListType string

const (
	ListTypeUnordered ListType = "unordered"
	ListTypeOrdered   ListType = "ordered"
)

func getContentDirectory() string {
	return fmt.Sprintf("blog-content/%s/content", env.MustEnvironmentName())
}

func getContentFileNameforID(id ID) string {
	return fmt.Sprintf("%s.json", id)
}

func GetContentForBlog(tx *sqlx.Tx, urlPath string, requireIsPostLive bool) ([]ContentNode, error) {
	metadata, err := getBlogPostMetadataByURLPath(tx, urlPath)
	if err != nil {
		return nil, err
	}
	if requireIsPostLive && metadata.Status != PostStatusLive {
		return nil, fmt.Errorf("Post is not live")
	}
	bytes, err := storage.RemoteStorage.Read(getContentDirectory(), getContentFileNameforID(metadata.ID))
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchKey") {
			return nil, nil
		}
		return nil, err
	}
	var content []ContentNode
	if err := json.Unmarshal(bytes, &content); err != nil {
		return nil, err
	}
	return content, nil
}

func UpsertContentForBlog(tx *sqlx.Tx, urlPath string, content []ContentNode) error {
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
	file, err := storage.NewFile(getContentFileNameforID(metadata.ID), bytes)
	if err != nil {
		return err
	}
	return storage.RemoteStorage.Write(getContentDirectory(), *file)
}

func verifyContent(content []ContentNode) error {
	var errs []string
	for idx, node := range content {
		bytes, err := json.Marshal(node.Body)
		if err != nil {
			errs = append(errs, fmt.Sprintf("Error marshalling node %d: %s", idx, err.Error()))
			continue
		}
		switch node.Type {
		case ContentNodeTypeHeading:
			var h Heading
			if err := json.Unmarshal(bytes, &h); err != nil {
				errs = append(errs, fmt.Sprintf("Node %d has type heading, but the body does not marshal correctly", idx))
			}
		case ContentNodeTypeParagraph:
			var p Paragraph
			if err := json.Unmarshal(bytes, &p); err != nil {
				errs = append(errs, fmt.Sprintf("Node %d has type paragraph, but the body does not marshal correctly", idx))
			}
		case ContentNodeTypeImage:
			var i Image
			if err := json.Unmarshal(bytes, &i); err != nil {
				errs = append(errs, fmt.Sprintf("Node %d has type image, but the body does not marshal correctly", idx))
			}
		case ContentNodeTypeLink:
			var l Link
			if err := json.Unmarshal(bytes, &l); err != nil {
				errs = append(errs, fmt.Sprintf("Node %d has type link, but the body does not marshal correctly", idx))
			}
		case ContentNodeTypeList:
			var l List
			if err := json.Unmarshal(bytes, &l); err != nil {
				errs = append(errs, fmt.Sprintf("Node %d has type list, but the body does not marshal correctly", idx))
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
