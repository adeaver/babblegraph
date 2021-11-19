package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Request struct {
	r *http.Request

	respStatus  *int
	respCookies []*http.Cookie
}

func (r *Request) GetJSONBody(v interface{}) error {
	bytes, err := ioutil.ReadAll(r.r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
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
