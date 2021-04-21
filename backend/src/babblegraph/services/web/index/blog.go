package index

import "net/http"

func ServeBlogPost(blogPostTemplateFileName string, blogPostURLPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handleUTMParameters(w, r)
		serveIndexTemplate(blogPostTemplateFileName, w, r)
	}
}
