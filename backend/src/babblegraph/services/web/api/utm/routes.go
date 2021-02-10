package utm

import (
	"babblegraph/services/web/router"
	"babblegraph/util/ptr"
	"encoding/json"
)

func RegisterRouteGroups() error {
	return router.RegisterRouteGroup(router.RouteGroup{
		Prefix: "utm",
		Routes: []router.Route{
			{
				Path:             "set_page_load_event_1",
				Handler:          handleSetPageLoadEvent,
				TrackEventWithID: ptr.String("page-load"),
			},
		},
	})
}

type setPageLoadEventRequest struct{}

type setPageLoadEventResponse struct{}

func handleSetPageLoadEvent(body []byte) (interface{}, error) {
	var req setPageLoadEventRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}
	return setPageLoadEventResponse{}, nil
}
