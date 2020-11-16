package scheduler

import (
	"runtime/debug"
	"time"

	cron "github.com/robfig/cron/v3"
)

func StartScheduler(errs chan error) error {
	usEastern, err := time.LoadLocation("America/New_York")
	if err != nil {
		return err
	}
	c := cron.New(cron.WithLocation(usEastern))
	c.AddFunc("30 2 * * *", makeRefetchSeedDomainJob(errs))
	c.Start()
	return nil
}

func makeRefetchSeedDomainJob(errs chan error) func() {
	return func() {
		defer func() {
			x := recover()
			if err, ok := x.(error); ok {
				errs <- err
				debug.PrintStack()
			}
		}()
		if err := RefetchSeedDomainsForNewContent(); err != nil {
			errs <- err
		}
	}
}
