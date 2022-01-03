package blog

import (
	"babblegraph/model/admin"
	"babblegraph/model/blog"
	"babblegraph/services/web/adminrouter/middleware"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/storage"
	"bytes"
	"io"

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
		}, {
			Path: "update_blog_content_1",
			Handler: middleware.WithPermission(
				admin.PermissionWriteBlog,
				updateBlogContent,
			),
		}, {
			Path: "get_blog_content_1",
			Handler: middleware.WithPermission(
				admin.PermissionWriteBlog,
				getBlogContent,
			),
		}, {
			Path: "upload_blog_image_1",
			Handler: middleware.WithPermission(
				admin.PermissionWriteBlog,
				uploadBlogImage,
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

type updateBlogContentRequest struct {
	URLPath string             `json:"url_path"`
	Content []blog.ContentNode `json:"content"`
}

type updateBlogContentResponse struct {
	Success bool `json:"success"`
}

func updateBlogContent(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req updateBlogContentRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	s3Storage := storage.NewS3StorageForEnvironment()
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		return blog.UpsertContentForBlog(tx, s3Storage, req.URLPath, req.Content)
	}); err != nil {
		return nil, err
	}
	return updateBlogContentResponse{
		Success: true,
	}, nil
}

type getBlogContentRequest struct {
	URLPath string `json:"url_path"`
}

type getBlogContentResponse struct {
	Content []blog.ContentNode `json:"content"`
}

func getBlogContent(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req getBlogContentRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	s3Storage := storage.NewS3StorageForEnvironment()
	var content []blog.ContentNode
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		content, err = blog.GetContentForBlog(tx, s3Storage, req.URLPath)
		return err
	}); err != nil {
		return nil, err
	}
	return getBlogContentResponse{
		Content: content,
	}, nil
}

type uploadBlogImageResponse struct {
}

func uploadBlogImage(adminID admin.ID, r *router.Request) (interface{}, error) {
	file, _, err := r.GetFile(nil)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}
	if err := database.WithTx(func(tx *sqlx.Tx) error {
	}); err != nil {
		return nil, err
	}
	return uploadBlogImageResponse{}, nil
}
