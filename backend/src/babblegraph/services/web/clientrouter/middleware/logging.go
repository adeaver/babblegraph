package middleware

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func WithBodyLogger(muxRouter func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		LogRequestWithBody(r)
		muxRouter(w, r)
	}
}

func WithoutBodyLogger(muxRouter func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		LogRequestWithoutBody(r)
		muxRouter(w, r)
	}
}

func LogRequestWithBody(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println(fmt.Sprintf("Error dumping request: %s", err.Error()))
		return
	}
	log.Println(string(requestDump))
}

func LogRequestWithoutBody(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, false)
	if err != nil {
		log.Println(fmt.Sprintf("Error dumping request: %s", err.Error()))
		return
	}
	log.Println(string(requestDump))
}
