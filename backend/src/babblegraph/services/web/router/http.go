package router

import (
	"babblegraph/util/deref"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
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

func (r *Request) GetJSONBody(v interface{}) error {
	bytes, err := ioutil.ReadAll(r.r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}

func (r *Request) GetFile(formFileFieldName *string) (multipart.File, *multipart.FileHeader, error) {
	return r.r.FormFile(deref.String(formFileFieldName, "file"))
}

func (r *Request) GetFormValue(fieldName string) string {
	return r.r.FormValue(fieldName)
}

func (r *Request) GetCookies() []*http.Cookie {
	return r.r.Cookies()
}

func (r *Request) RespondWithCookie(c *http.Cookie) {
	r.respCookies = append(r.respCookies, c)
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
