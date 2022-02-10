package router

import (
	"babblegraph/util/bglog"
	"context"
)

type Context struct {
	ctx    context.Context
	logger *bglog.Logger
}

func (c Context) Debugf(format string, args ...interface{}) {
	c.logger.Debugf(format, args...)
}

func (c Context) Infof(format string, args ...interface{}) {
	c.logger.Infof(format, args...)
}

func (c Context) Warnf(format string, args ...interface{}) {
	c.logger.Warnf(format, args...)
}

func (c Context) Errorf(format string, args ...interface{}) {
	c.logger.Errorf(format, args...)
}
