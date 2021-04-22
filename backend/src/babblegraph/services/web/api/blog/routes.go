package blog

import (
	"babblegraph/model/blogposts"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"encoding/json"
	"log"

	"github.com/jmoiron/sqlx"
)

func RegisterRouteGroups() error {
	return router.RegisterRouteGroup(router.RouteGroup{
		Prefix: "blog",
		Routes: []router.Route{
			{
				Path:    "get_all_blog_posts_paginated_1",
				Handler: handleGetAllBlogPostsPaginated,
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
	log.Println("Reach here")
	return getAllBlogPostsPaginatedResponse{
		BlogPosts: blogPosts,
	}, nil
}
