package contentingestion

import (
	"babblegraph/model/links2"
	"babblegraph/util/ctx"
)

func processWebsiteHTML1Link(c ctx.LogContext, link links2.Link) error {
	c.Infof("Would process link")
	return nil
}
