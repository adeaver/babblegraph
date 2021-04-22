package index

import (
	"babblegraph/model/blogposts"
	"babblegraph/model/utm"
	"babblegraph/util/database"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/jmoiron/sqlx"
)

func ServeBlogPost(blogPostTemplateFileName string, blogPostURLPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var b *blogposts.BlogPost
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			var err error
			b, err = blogposts.LookupBlogByURLPath(tx, blogPostURLPath)
			if err != nil {
				return err
			}
			if b == nil {
				return nil
			}
			return blogposts.CaptureBlogPostView(tx, b.ID)
		}); err != nil {
			sentry.CaptureException(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if b == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var utmParameters *utm.Parameters
		utmParameters = utm.GetParametersForRequest(r)
		if utmParameters == nil {
			utmParameters = &utm.Parameters{
				Source:     utm.Source("blog").Ptr(),
				CampaignID: utm.CampaignID(fmt.Sprintf("blog-post-%s", b.TrackingTag)).Ptr(),
			}
		}
		handleUTMParameters(*utmParameters, w, r)
		serveIndexTemplate(blogPostTemplateFileName, *b, w, r)
	}
}

func addBlogPostUTMParameters(r *http.Request) {

}
