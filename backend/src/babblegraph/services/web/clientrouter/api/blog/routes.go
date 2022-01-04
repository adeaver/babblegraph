package blog

import (
	"babblegraph/model/blog"
	"babblegraph/services/web/clientrouter/api"
	"babblegraph/util/database"
	"encoding/json"

	"github.com/jmoiron/sqlx"
)

func RegisterRouteGroups() error {
	return api.RegisterRouteGroup(api.RouteGroup{
		Prefix: "blog",
		Routes: []api.Route{
			{
				Path:    "get_blog_metadata_1",
				Handler: handleGetBlogMetadata,
			}, {
				Path:    "get_blog_content_1",
				Handler: handleGetBlogContent,
			},
		},
	})
}

type getBlogMetadataRequest struct {
	URLPath string `json:"url_path"`
}

type getBlogMetadataResponse struct {
	BlogPost blog.BlogPostMetadata `json:"blog_post"`
}

func handleGetBlogMetadata(body []byte) (interface{}, error) {
	var req getBlogMetadataRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	var blogPostMetadata *blog.BlogPostMetadata
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		blogPostMetadata, err = blog.GetClientBlogPostMetadataForURLPath(tx, req.URLPath)
		return err
	}); err != nil {
		return nil, err
	}
	return getBlogMetadataResponse{
		BlogPost: *blogPostMetadata,
	}, nil
}

type getBlogContentRequest struct {
	URLPath string `json:"url_path"`
}

type getBlogContentResponse struct {
	Content []blog.ContentNode `json:"content"`
}

func handleGetBlogContent(body []byte) (interface{}, error) {
	var req getBlogContentRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	var content []blog.ContentNode
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		content, err = blog.GetContentForBlog(tx, req.URLPath, true)
		return err
	}); err != nil {
		return nil, err
	}
	return getBlogContentResponse{
		Content: content,
	}, nil
}
