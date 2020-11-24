package opengraph

import "time"

type BasicMetadata struct {
	Title           *string
	Type            *string
	ImageURL        *string
	URL             *string
	Description     *string
	PublicationTime *time.Time
}

func GetBasicMetadata(metadata map[string]string) BasicMetadata {
	return BasicMetadata{
		Title:           findTagOrNil(TitleTag, metadata),
		Type:            findTagOrNil(TypeTag, metadata),
		ImageURL:        findTagOrNil(ImageURLTag, metadata),
		URL:             findTagOrNil(URLTag, metadata),
		Description:     findTagOrNil(DescriptionTag, metadata),
		PublicationTime: lookupPublicationTime(metadata),
	}
}

func findTagOrNil(tag Tag, metadata map[string]string) *string {
	if val, ok := metadata[tag.Str()]; ok {
		return &val
	}
	return nil
}

func lookupPublicationTime(metadata map[string]string) *time.Time {
	if strTime, ok := metadata[PublicationTimeTag.Str()]; ok {
		// According to the opengraph protocol (https://ogp.me/)
		// All times are in ISO 8601, which is RFC 3339
		t, err := time.Parse(time.RFC3339, strTime)
		if err != nil {
			// I think we're relying on random people
			// to input time potentially, so if it's malformed
			// we should just accept defeat.
			return nil
		}
		return &t
	}
	return nil
}
