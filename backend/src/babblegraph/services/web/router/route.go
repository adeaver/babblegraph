package router

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

type Route struct {
	Path          string
	Handler       RouteHandler
	UseCaptcha    bool
	ShouldLogBody bool
}

type RouteHandler func(reqBody []byte) (_resp interface{}, _err error)

func makeMuxRouter(processRequest RouteHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(fmt.Sprintf("ERROR: %s", err.Error()))
			writeErrorJSONResponse(w, errorResponse{
				Message: "Request is not valid",
			})
			return
		}
		resp, err := processRequest(body)
		if err != nil {
			log.Println(fmt.Sprintf("ERROR: %s", err.Error()))
			writeErrorJSONResponse(w, errorResponse{
				Message: "Error processing request",
			})
			return
		}
		w.WriteHeader(http.StatusOK)
		writeJSONResponse(w, resp)
	}
}

func withCaptcha(muxRouter func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Captcha")
		muxRouter(w, r)
	}
}

func withBodyLogger(muxRouter func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		LogRequestWithBody(r)
		muxRouter(w, r)
	}
}

func withoutBodyLogger(muxRouter func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
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
