package routermiddleware

import "babblegraph/services/web/router"

func WithRequestBodyLogger(handler router.RequestHandler) router.RequestHandler {
	return func(r *router.Request) (interface{}, error) {
		r.LogRequest(true)
		return handler(r)
	}
}

func WithNoBodyRequestLogger(handler router.RequestHandler) router.RequestHandler {
	return func(r *router.Request) (interface{}, error) {
		r.LogRequest(false)
		return handler(r)
	}
}
