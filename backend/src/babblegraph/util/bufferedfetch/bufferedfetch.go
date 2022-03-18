package bufferedfetch

import (
	"babblegraph/util/ctx"
	"fmt"
	"log"
	"reflect"
)

type refillFunc func() (interface{}, error)

type fetchState struct {
	i        []interface{}
	refillFn refillFunc
}

var bufferedState = map[string]*fetchState{}

func Register(c ctx.LogContext, key string, fn refillFunc) error {
	if _, exists := bufferedState[key]; exists {
		c.Infof("Cannot register bufferedfetch for key %s: already exists", key)
		return nil
	}
	bufferedState[key] = &fetchState{
		refillFn: fn,
	}
	return nil
}

func WithNextBufferedValue(key string, fn func(interface{}) error) error {
	fs, exists := bufferedState[key]
	if !exists {
		return fmt.Errorf("Cannot get next value for bufferedfetch with key %s: does not exist", key)
	}
	if len(fs.i) == 0 {
		log.Println(fmt.Sprintf("Bufferedfetch refilling key: %s", key))
		refill, err := fs.refillFn()
		if err != nil {
			return fmt.Errorf("Error refilling bufferedfetch with key %s: %s", key, err.Error())
		}
		switch reflect.TypeOf(refill).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(refill)
			for i := 0; i < s.Len(); i++ {
				fs.i = append(fs.i, s.Index(i).Interface())
			}
		default:
			return fmt.Errorf("Error refilling bufferedfetch with key %s: refill func did not return a list", key)
		}
	}
	if len(fs.i) == 0 {
		return nil
	}
	nextVal := fs.i[0]
	fs.i = append([]interface{}{}, fs.i[1:]...)
	return fn(nextVal)
}

func ForceRefill(c ctx.LogContext, key string) error {
	fs, exists := bufferedState[key]
	if !exists {
		return fmt.Errorf("Cannot get next value for bufferedfetch with key %s: does not exist", key)
	}
	c.Infof("Bufferedfetch refilling key %s on force", key)
	refill, err := fs.refillFn()
	if err != nil {
		return fmt.Errorf("Error refilling bufferedfetch with key %s: %s", key, err.Error())
	}
	switch reflect.TypeOf(refill).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(refill)
		for i := 0; i < s.Len(); i++ {
			fs.i = append(fs.i, s.Index(i).Interface())
		}
	default:
		return fmt.Errorf("Error refilling bufferedfetch with key %s: refill func did not return a list", key)
	}
	return nil
}
