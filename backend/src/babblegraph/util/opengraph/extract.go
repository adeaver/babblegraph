package opengraph

type BasicMetadata struct {
	Title    *string
	Type     *string
	ImageURL *string
	URL      *string
}

func GetBasicMetadata(metadata map[string]string) BasicMetadata {
	return BasicMetadata{
		Title:    findTagOrNil(TitleTag, metadata),
		Type:     findTagOrNil(TypeTag, metadata),
		ImageURL: findTagOrNil(ImageURLTag, metadata),
		URL:      findTagOrNil(URLTag, metadata),
	}
}

func findTagOrNil(tag Tag, metadata map[string]string) *string {
	if val, ok := metadata[tag.Str()]; ok {
		return &val
	}
	return nil
}
