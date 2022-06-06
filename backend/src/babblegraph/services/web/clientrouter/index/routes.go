package index

import (
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/router"
)

func GetRoutes() []router.IndexRoute {
	routes := []router.IndexRoute{
		{
			Path: router.IndexPath{
				Text: "/vfile/{fileName}",
			},
			Handler: routermiddleware.WithNoBodyRequestLogger(
				handleVirtualFile,
			),
		}, {
			Path: router.IndexPath{
				Text: "/verify/{token}",
			},
			Handler: routermiddleware.WithNoBodyRequestLogger(
				handleVerification,
			),
		}, {
			Path: router.IndexPath{
				Text: "/link/{token}",
			},
			Handler: routermiddleware.WithNoBodyRequestLogger(
				handleAdClick,
			),
		},
	}
	routes = append(routes, getPromotionCodeRoutes()...)
	return routes
}
