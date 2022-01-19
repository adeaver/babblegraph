package index

import (
	"babblegraph/model/blog"
	"babblegraph/model/utm"
	"babblegraph/util/async"
	"babblegraph/util/ctx"
	"babblegraph/util/database"
	"babblegraph/util/ptr"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func HandleServeBlogPost(staticFileDirName string) func(w http.ResponseWriter, r *http.Request) {
	c := ctx.GetDefaultLogContext()
	return func(w http.ResponseWriter, r *http.Request) {
		routeVars := mux.Vars(r)
		blogURLPath, ok := routeVars["blog_path"]
		if !ok {
			// Display general blog page
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var blogPostMetadata *blog.BlogPostMetadata
		if err := database.WithTx(func(tx *sqlx.Tx) error {
			var err error
			blogPostMetadata, err = blog.GetClientBlogPostMetadataForURLPath(tx, blogURLPath)
			return err
		}); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		utmTrackingID, utmParameters, err := handleUTMParametersForBlogPost(blogPostMetadata, r)
		switch {
		case err != nil:
			c.Errorf("Error getting tracking ID: %s", err.Error())
		case utmTrackingID != nil:
			http.SetCookie(w, &http.Cookie{
				Name:  utm.UTMTrackingIDCookieName,
				Value: utmTrackingID.Str(),
			})
		default:
			c.Infof("No tracking ID, but also no error")
		}
		errs := make(chan error)
		async.WithContext(errs, "serve-blog", func(c async.Context) {
			if err := database.WithTx(func(tx *sqlx.Tx) error {
				if err := blog.RegisterBlogView(tx, blogPostMetadata.ID, utmTrackingID); err != nil {
					return err
				}
				if utmTrackingID == nil {
					return nil
				}
				return utm.RegisterUTMPageHit(tx, *utmTrackingID, *utmParameters)
			}); err != nil {
				c.Errorf("Error handling UTM Page hit: %s", err.Error())
			}
		}).Start()
		w.Header().Add("Content-Type", "text/html")
		var tmpl *template.Template
		tmpl, err = template.New("blog_index.html").ParseFiles(fmt.Sprintf("%s/blog_index.html", staticFileDirName))
		if err != nil {
			c.Errorf("Error loading blog index page: %s", err.Error())
			http.Error(w, http.StatusText(500), 500)
			return
		}
		err = tmpl.Execute(w, *blogPostMetadata)
		if err != nil {
			c.Errorf("Error executing blog index template: %s", err.Error())
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}
}

func handleUTMParametersForBlogPost(blogPostMetadata *blog.BlogPostMetadata, r *http.Request) (*utm.TrackingID, *utm.Parameters, error) {
	utmParameters := utm.GetParametersForRequest(r)
	if utmParameters == nil {
		utmParameters = &utm.Parameters{
			Source:     utm.Source("blog").Ptr(),
			CampaignID: utm.CampaignID(string(blogPostMetadata.ID)).Ptr(),
			Medium:     utm.Medium("organic").Ptr(),
			URLPath:    ptr.String(blogPostMetadata.URLPath),
		}
	}
	trackingID, err := utm.GetTrackingIDForRequest(r)
	return trackingID, utmParameters, err
}
