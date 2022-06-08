package router

import (
	"babblegraph/util/deref"
	"babblegraph/util/ptr"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gorilla/mux"
)

type Request struct {
	c Context
	r *http.Request

	respStatus  *int
	respCookies []*http.Cookie
}

func (r *Request) LogRequest(includeBody bool) {
	requestDump, err := httputil.DumpRequest(r.r, includeBody)
	if err != nil {
		r.Warnf("Error dumping request: %s", err.Error())
		return
	}
	r.Infof(string(requestDump))
}

func (r *Request) GetBodyAsBytes() ([]byte, error) {
	bytes, err := ioutil.ReadAll(r.r.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (r *Request) GetJSONBody(v interface{}) error {
	bytes, err := r.GetBodyAsBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

func (r *Request) GetFile(formFileFieldName *string) (multipart.File, *multipart.FileHeader, error) {
	return r.r.FormFile(deref.String(formFileFieldName, "file"))
}

func (r *Request) GetRouteVar(varName string) (*string, error) {
	routeVars := mux.Vars(r.r)
	routeVar, ok := routeVars[varName]
	if !ok {
		return nil, fmt.Errorf("Route var %s not found", varName)
	}
	return ptr.String(routeVar), nil
}

func (r *Request) GetQueryParam(varName string) *string {
	param := r.r.URL.Query().Get(varName)
	if param != "" {
		return ptr.String(param)
	}
	return nil
}

func (r *Request) GetFormValue(fieldName string) string {
	return r.r.FormValue(fieldName)
}

func (r *Request) GetHeader(headerName string) string {
	return r.r.Header.Get(headerName)
}

func (r *Request) GetCookies() []*http.Cookie {
	return r.r.Cookies()
}

func (r *Request) RespondWithCookie(c *http.Cookie) {
	r.respCookies = append(r.respCookies, c)
}

func (r *Request) RemoveCookieByName(cookieName string) {
	r.respCookies = append(r.respCookies, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(-5 * time.Minute),
	})
}

func (r *Request) RespondWithStatus(status int) {
	r.respStatus = &status
}

func (r *Request) Debugf(format string, args ...interface{}) {
	r.c.Debugf(format, args...)
}

func (r *Request) Infof(format string, args ...interface{}) {
	r.c.Infof(format, args...)
}

func (r *Request) Warnf(format string, args ...interface{}) {
	r.c.Warnf(format, args...)
}

func (r *Request) Errorf(format string, args ...interface{}) {
	r.c.Errorf(format, args...)
}
