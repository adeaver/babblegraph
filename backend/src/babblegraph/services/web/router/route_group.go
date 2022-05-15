package router

import (
	"babblegraph/util/bglog"
	"babblegraph/util/random"
	"context"
	"html/template"
	"net/http"
)

type RouteGroup struct {
	Prefix string
	Routes []Route
}

type Route struct {
	Path    string
	Handler RequestHandler
}

type RequestHandler func(r *Request) (interface{}, error)

func (r Route) makeMuxRoute() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		contextKey := random.MustMakeRandomString(12)
		ctx := Context{
			ctx:    context.Background(),
			logger: bglog.NewLoggerForContext(r.Path, contextKey, 4),
		}
		wrappedRequest := &Request{
			c: ctx,
			r: req,
		}
		status := http.StatusOK
		resp, err := r.Handler(wrappedRequest)
		switch {
		case err != nil:
			ctx.Errorf("Got error processing request: %s", err.Error())
			writeErrorJSONResponse(w, errorResponse{
				Message: "Error processing request",
			})
			return
		case wrappedRequest.respStatus != nil:
			status = *wrappedRequest.respStatus
		}
		if len(wrappedRequest.respCookies) != 0 {
			for _, cookie := range wrappedRequest.respCookies {
				http.SetCookie(w, cookie)
			}
		}
		w.WriteHeader(status)
		writeJSONResponse(w, resp)
	}
}

type IndexRoute struct {
	Path    IndexPath
	Handler RequestHandler
}

type IndexPath struct {
	Text        string
	UseAsPrefix bool
}

type IndexResponse struct {
	ContentType  *ContentType
	FileTemplate *template.Template
	TemplateData interface{}
}

type ContentType string

const (
	ContentTypeTextHTML ContentType = "text/html"
)

func (c ContentType) Ptr() *ContentType {
	return &c
}

func (c ContentType) Str() string {
	return string(c)
}

func (i IndexRoute) makeMuxRoute() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		contextKey := random.MustMakeRandomString(12)
		c := Context{
			ctx:    context.Background(),
			logger: bglog.NewLoggerForContext(i.Path.Text, contextKey, 4),
		}
		wrappedRequest := &Request{
			c: c,
			r: req,
		}
		resp, err := i.Handler(wrappedRequest)
		if err != nil {
			c.Errorf("Error processing request: %s", err.Error())
			http.Error(w, http.StatusText(500), 500)
			return
		}
		indexResponse, ok := resp.(*IndexResponse)
		if ok {
			contentType := ContentTypeTextHTML
			if indexResponse.ContentType != nil {
				contentType = *indexResponse.ContentType
			}
			if len(wrappedRequest.respCookies) != 0 {
				for _, cookie := range wrappedRequest.respCookies {
					http.SetCookie(w, cookie)
				}
			}
			w.Header().Add("Content-Type", contentType.Str())
			err = indexResponse.FileTemplate.Execute(w, indexResponse.TemplateData)
			if err != nil {
				c.Errorf("Error executing template for path %s: %s", i.Path.Text, err.Error())
				http.Error(w, http.StatusText(500), 500)
				return
			}
			return
		}
		redirectURL, ok := resp.(*string)
		if ok {
			if len(wrappedRequest.respCookies) != 0 {
				for _, cookie := range wrappedRequest.respCookies {
					http.SetCookie(w, cookie)
				}
			}
			http.Redirect(w, req, *redirectURL, http.StatusFound)
			return
		}
		c.Errorf("Did not get implementation for path: %s", i.Path.Text)
		http.Error(w, http.StatusText(500), 500)
	}
}
