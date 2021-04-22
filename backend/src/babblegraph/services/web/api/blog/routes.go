package blog

import (
	"babblegraph/model/blogposts"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/storage"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func RegisterRouteGroups() error {
	return router.RegisterRouteGroup(router.RouteGroup{
		Prefix: "blog",
		Routes: []router.Route{
			{
				Path:    "get_all_blog_posts_paginated_1",
				Handler: handleGetAllBlogPostsPaginated,
			}, {
				Path:    "get_blog_post_data_1",
				Handler: handleGetBlogPostData,
			},
		},
	})
}

type getAllBlogPostsPaginatedRequest struct {
	PageIndex int `json:"page_index"`
}

type getAllBlogPostsPaginatedResponse struct {
	BlogPosts []blogposts.BlogPost `json:"blog_posts"`
}

func handleGetAllBlogPostsPaginated(body []byte) (interface{}, error) {
	var req getAllBlogPostsPaginatedRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	var blogPosts []blogposts.BlogPost
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		blogPosts, err = blogposts.GetAllBlogsPaginated(tx, req.PageIndex)
		return err
	}); err != nil {
		return nil, err
	}
	return getAllBlogPostsPaginatedResponse{
		BlogPosts: blogPosts,
	}, nil
}

type getBlogPostDataRequest struct {
	URLPath string `json:"url_path"`
}

type getBlogPostDataResponse struct {
	Metadata blogposts.BlogPost `json:"metadata"`
	Content  string             `json:"content"`
}

func handleGetBlogPostData(body []byte) (interface{}, error) {
	var req getBlogPostDataRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	var blogPost *blogposts.BlogPost
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		blogPost, err = blogposts.LookupBlogByURLPath(tx, req.URLPath)
		return err
	}); err != nil {
		return nil, err
	}
	if blogPost == nil {
		return nil, fmt.Errorf("No blog post found for url path: %s", req.URLPath)
	}
	content, err := getBlogPostContent(blogPost.ContentURL)
	if err != nil {
		return nil, err
	}
	return getBlogPostDataResponse{
		Metadata: *blogPost,
		Content:  *content,
	}, nil
}

func getBlogPostContent(contentURL string) (*string, error) {
	return storage.NewS3StorageForEnvironment().GetData("prod-spaces-1/blog/content", contentURL)
}
