package blog

import (
	"babblegraph/model/admin"
	"babblegraph/model/blog"
	"babblegraph/services/web/adminrouter/middleware"
	"babblegraph/services/web/router"
	"babblegraph/util/database"

	"github.com/jmoiron/sqlx"
)

var Routes = router.RouteGroup{
	Prefix: "blog",
	Routes: []router.Route{
		{
			Path: "get_all_blog_post_metadata_1",
			Handler: middleware.WithPermission(
				admin.PermissionWriteBlog,
				getAllBlogPostMetadata,
			),
		}, {
			Path: "get_blog_post_metadata_by_url_path_1",
			Handler: middleware.WithPermission(
				admin.PermissionWriteBlog,
				getBlogPostMetadataByURLPath,
			),
		}, {
			Path: "add_blog_post_metadata_1",
			Handler: middleware.WithPermission(
				admin.PermissionWriteBlog,
				addBlogPostMetadata,
			),
		}, {
			Path: "update_blog_post_metadata_1",
			Handler: middleware.WithPermission(
				admin.PermissionWriteBlog,
				updateBlogPostMetadata,
			),
		}, {
			Path: "update_blog_post_status_1",
			Handler: middleware.WithPermission(
				admin.PermissionPublishBlog,
				updateBlogPostStatus,
			),
		},
	},
}

type getAllBlogPostMetadataResponse struct {
	AllBlogPosts []blog.BlogPostMetadata `json:"all_blog_posts"`
}

func getAllBlogPostMetadata(adminID admin.ID, r *router.Request) (interface{}, error) {
	var blogPosts []blog.BlogPostMetadata
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		blogPosts, err = blog.GetAllBlogPostMetadata(tx)
		return err
	}); err != nil {
		return nil, err
	}
	return getAllBlogPostMetadataResponse{
		AllBlogPosts: blogPosts,
	}, nil
}

type getBlogPostMetadataByURLPathRequest struct {
	URLPath string `json:"url_path"`
}

type getBlogPostMetadataByURLPathResponse struct {
	BlogPost blog.BlogPostMetadata `json:"blog_post"`
}

func getBlogPostMetadataByURLPath(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req getBlogPostMetadataByURLPathRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var blogPost *blog.BlogPostMetadata
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		blogPost, err = blog.GetBlogPostMetadataByURLPath(tx, req.URLPath)
		return err
	}); err != nil {
		return nil, err
	}
	return getBlogPostMetadataByURLPathResponse{
		BlogPost: *blogPost,
	}, nil
}

type addBlogPostMetadataRequest struct {
	URLPath     string `json:"url_path"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorName  string `json:"author_name"`
}

type addBlogPostMetadataResponse struct {
	Success bool `json:"success"`
}

func addBlogPostMetadata(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req addBlogPostMetadataRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return blog.AddBlogPostMetadata(tx, blog.AddBlogPostMetadataInput{
			Title:       req.Title,
			URLPath:     req.URLPath,
			AuthorName:  req.AuthorName,
			Description: req.Description,
		})
	}); err != nil {
		return nil, err
	}
	return addBlogPostMetadataResponse{
		Success: true,
	}, nil
}

type updateBlogPostMetadataRequest struct {
	URLPath       string  `json:"url_path"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	HeroImagePath *string `json:"hero_image_path,omitempty"`
	AuthorName    string  `json:"author_name"`
}

type updateBlogPostMetadataResponse struct {
	Success bool `json:"success"`
}

func updateBlogPostMetadata(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req updateBlogPostMetadataRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return blog.UpdateBlogPostMetadata(tx, blog.UpdateBlogPostMetadataInput{
			Title:         req.Title,
			URLPath:       req.URLPath,
			AuthorName:    req.AuthorName,
			Description:   req.Description,
			HeroImagePath: req.HeroImagePath,
		})
	}); err != nil {
		return nil, err
	}
	return updateBlogPostMetadataResponse{
		Success: true,
	}, nil
}

type updateBlogPostStatusRequest struct {
	URLPath string          `json:"url_path"`
	Status  blog.PostStatus `json:"status"`
}

type updateBlogPostStatusResponse struct {
	Success bool `json:"success"`
}

func updateBlogPostStatus(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req updateBlogPostStatusRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return blog.UpdateBlogPostStatus(tx, req.URLPath, req.Status)
	}); err != nil {
		return nil, err
	}
	return updateBlogPostStatusResponse{
		Success: true,
	}, nil
}
