package bufferedfetch

func Register(key string, fn func() (interface{}, error)) error {
	return nil
}

func WithNextBufferedValue(key string, func (interface{}) error) error {
    return nil
}
