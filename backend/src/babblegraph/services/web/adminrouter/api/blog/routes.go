package blog

import (
	"babblegraph/model/admin"
	"babblegraph/model/blog"
	"babblegraph/services/web/adminrouter/middleware"
	"babblegraph/services/web/router"
	"babblegraph/util/database"
	"babblegraph/util/ptr"
	"babblegraph/util/storage"
	"bytes"
	"fmt"
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
		}, {
			Path: "get_blog_post_view_metrics_1",
			Handler: middleware.WithPermission(
				admin.PermissionWriteBlog,
				getBlogPostViewMetrics,
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
	URLPath     string `json:"url_path"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorName  string `json:"author_name"`
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
			Title:       req.Title,
			URLPath:     req.URLPath,
			AuthorName:  req.AuthorName,
			Description: req.Description,
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
	var content []blog.ContentNode
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		content, err = blog.GetContentForBlog(tx, req.URLPath, false)
		return err
	}); err != nil {
		return nil, err
	}
	return getBlogContentResponse{
		Content: content,
	}, nil
}

type uploadBlogImageResponse struct {
	ImagePath string `json:"image_path"`
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
	blogURLPath := r.GetFormValue("url_path")
	var caption *string
	if captionValue := r.GetFormValue("caption"); len(captionValue) != 0 {
		caption = ptr.String(captionValue)
	}
	fileName := r.GetFormValue("file_name")
	storageFile, err := storage.NewFile(fileName, buf.Bytes())
	if err != nil {
		return nil, err
	}
	isHeroImage := r.GetFormValue("is_hero_image") == "true"
	storageFile.AssignAccessControlLevel(storage.AccessControlPublicReadOnly)
	var path string
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		blogPostMetadata, err := blog.GetBlogPostMetadataByURLPath(tx, blogURLPath)
		if err != nil {
			return err
		}
		imageDirectory := blog.MakeImageDirectory(blogPostMetadata.ID)
		path = fmt.Sprintf("%s/%s", imageDirectory, fileName)
		if err := blog.InsertBlogImageMetadata(tx, blog.InsertBlogImageMetadataInput{
			BlogID:      blogPostMetadata.ID,
			Path:        path,
			FileName:    fileName,
			AltText:     r.GetFormValue("alt_text"),
			Caption:     caption,
			IsHeroImage: isHeroImage,
		}); err != nil {
			return err
		}
		return storage.RemoteStorage.Write(imageDirectory, *storageFile)
	}); err != nil {
		return nil, err
	}
	return uploadBlogImageResponse{
		ImagePath: path,
	}, nil
}

type getBlogPostViewMetricsRequest struct {
	URLPath string `json:"url_path"`
}

type getBlogPostViewMetricsResponse struct {
	ViewMetrics blog.BlogPostViewMetrics `json:"view_metrics"`
}

func getBlogPostViewMetrics(adminID admin.ID, r *router.Request) (interface{}, error) {
	var req getBlogPostViewMetricsRequest
	if err := r.GetJSONBody(&req); err != nil {
		return nil, err
	}
	var viewMetrics *blog.BlogPostViewMetrics
	if err := database.WithTx(func(tx *sqlx.Tx) error {
		var err error
		viewMetrics, err = blog.GetBlogPostViewMetrics(tx, req.URLPath)
		return err
	}); err != nil {
		return nil, err
	}
	return getBlogPostViewMetricsResponse{
		ViewMetrics: *viewMetrics,
	}, nil
}
