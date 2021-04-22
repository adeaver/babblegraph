package index

import (
	"babblegraph/model/blogposts"
	"net/http"
)

func ServeBlogPost(blogPostTemplateFileName string, blogPostURLPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handleUTMParameters(w, r)
		b := blogposts.BlogPost{
			Title:            "Sample Blog Post",
			Description:      "Sample Description",
			HeroImageURL:     "sample-hero-image.jpg",
			HeroImageAltText: "Sample image for a blog",
			ContentURL:       "sample-blog.json",
		}
		serveIndexTemplate(blogPostTemplateFileName, b, w, r)
	}
}
