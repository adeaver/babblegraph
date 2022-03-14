package index

import (
	"babblegraph/services/web/clientrouter/routermiddleware"
	"babblegraph/services/web/router"
)

var Routes = []router.IndexRoute{
	{
		Path: router.IndexPath{
			Text: "/article/{token}",
		},
		Handler: routermiddleware.WithNoBodyRequestLogger(
			handleArticleRoute,
		),
	}, {
		Path: router.IndexPath{
			Text: "/paywall-report/{token}",
		},
		Handler: routermiddleware.WithNoBodyRequestLogger(
			handlePaywallReport,
		),
	}, {
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
	},
}
