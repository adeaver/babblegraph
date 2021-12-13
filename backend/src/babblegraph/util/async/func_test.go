package async

import (
	"fmt"
	"testing"
)

func TestFuncRuns(t *testing.T) {
	i := 0
	errs := make(chan error)
	WithContext(errs, "test", func(c Context) {
		i++
		errs <- fmt.Errorf("Done")
	}).Start()
	select {
	case err := <-errs:
		if err.Error() != "Done" {
			t.Errorf("Expected error Done, but got %s", err.Error())
		}
		if i != 1 {
			t.Errorf("Expected i to be 1, but got %d", i)
		}
	}
}

func TestFuncWithPanic(t *testing.T) {
	i := 0
	errs := make(chan error)
	WithContext(errs, "test", func(c Context) {
		c.Infof("Starting test")
		i++
		panic("This is a panic")
	}).Start()
	select {
	case _ = <-errs:
		if i != 1 {
			t.Errorf("Expected i to be 1, but got %d", i)
		}
	}
}
